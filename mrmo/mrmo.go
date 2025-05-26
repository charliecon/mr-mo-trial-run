package mrmo

import (
	"context"
	mockDynamo "github.com/charliecon/mr-mo-trial-run/mock-dynamo"
	orgManager "github.com/charliecon/mr-mo-trial-run/mrmo/org_manager"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	resourceExporter "github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/resource_exporter"
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/util"
	"log"
)

type MrMo struct {
	ResourceType   string
	Id             string
	ResourceData   *schema.ResourceData
	SchemaResource *schema.Resource
	ProviderMeta   any
	OrgManager     *orgManager.OrgManager
	Exporter       *resourceExporter.ResourceExporter

	// determined after export
	ResourcePath  string
	ResourceLabel string
}

type Message struct {
	ResourceType string
	EntityId     string
	IsDelete     bool
}

// ProcessMessage handles the processing of resource management operations based on incoming messages.
// It coordinates the export, configuration, and application of resources across target organizations.
//
// Parameters:
//   - ctx: Context for the operation
//   - message: Message struct containing:
//   - ResourceType: The type of resource being processed
//   - EntityId: The identifier of the source entity to process
//   - IsDelete: Flag indicating if this is a delete operation
//   - om: OrgManager instance for handling organization-specific operations
//
// Returns:
//   - diag.Diagnostics: Collection of diagnostic messages and errors encountered during processing
//
// The function performs the following sequence:
//  1. Initializes a new MrMo instance for the specified resource
//  2. For delete operations:
//     * Directly applies deletion across target organizations
//  3. For create/update operations:
//     * Exports the current resource configuration
//     * Parses the resource path from the configuration
//     * Appends necessary output blocks to the configuration (these are used to retrieve the target resource ID after apply)
//     * Applies the configuration across target organizations
func ProcessMessage(ctx context.Context, message Message, credentialsFilePath string) (diags diag.Diagnostics) {
	mrMo, err := newMrMo(message.ResourceType, credentialsFilePath, message.EntityId)
	if err != nil {
		return diag.FromErr(err)
	}

	if message.IsDelete {
		return mrMo.applyResourceConfigToTargetOrgs(nil, true)
	}

	resourceConfig, exportDiags := mrMo.exportConfig(ctx, message.EntityId, message.ResourceType)
	diags = append(diags, exportDiags...)
	if diags.HasError() {
		return diags
	}

	resourceLabel, err := parseResourceLabelFromConfig(resourceConfig, message.ResourceType)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	mrMo.ResourceLabel = resourceLabel
	mrMo.ResourcePath = message.ResourceType + "." + resourceLabel

	resourceConfig = appendOutputBlockToConfig(resourceConfig, mrMo.ResourcePath, message.EntityId)

	diags = append(diags, mrMo.applyResourceConfigToTargetOrgs(resourceConfig, false)...)
	return diags
}

// applyResourceConfigToTargetOrgs applies the provided resource configuration across all target organizations,
// handling credential management and GUID resolution for each target.
//
// Parameters:
//   - resourceConfig: JsonMap containing the exported configuration to be applied
//   - delete: Boolean flag indicating if this is a delete operation
//
// Returns:
//   - diag.Diagnostics: Collection of diagnostic messages and errors encountered during execution
//
// The function performs the following operations for each target organization:
//  1. Preserves original client credentials and restores them upon completion (the original client credentials are
//     restored via deferred function regardless of success or failure.)
//  2. Sets target organization-specific credentials
//  3. Resolves GUIDs in the resource configuration to match the target organization
//  4. Updates the Terraform configuration file in S3 for the target organization
//  5. Runs a targeted tofu apply
//  6. Updates the global mapping table according (add new mapping for create, delete for delete, etc.)
//
// If any operation fails for a target organization, the function returns immediately with error diagnostics.
// The original client credentials are restored via deferred function regardless of success or failure.
func (m *MrMo) applyResourceConfigToTargetOrgs(resourceConfig util.JsonMap, delete bool) (diags diag.Diagnostics) {
	originalClientId, originalClientSecret, originalRegion := orgManager.GetClientCredsEnvVars()
	defer func() {
		// restore client cred env vars
		err := orgManager.SetClientCredEnvVars(originalClientId, originalClientSecret, originalRegion)
		if err != nil {
			log.Printf("failed to restore client creds. Error: %s", err.Error())
		}
	}()

	for _, target := range m.OrgManager.Targets {
		// Set target org client credentials
		err := target.SetTargetOrgCredentials()
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}

		resourceConfigCopy := make(util.JsonMap)
		for k, v := range resourceConfig {
			resourceConfigCopy[k] = v
		}
		// Resolve GUIDs in resourceConfig to target org GUIDs using the mapping table
		if !delete {
			resourceConfigCopy, err = m.resolveResourceConfigDependencies(resourceConfigCopy, target)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		fm := newFileManager(target.OrgId, m.Id)

		// Update the tf file in s3 for the current target org
		diags = append(diags, fm.updateTargetTfConfig(resourceConfigCopy, delete)...)
		if diags.HasError() {
			return diags
		}

		// Run targeted apply
		targetResourceId, applyDiags := applyWithOpenTofu(fm.targetConfigDir, m.Id, m.ResourcePath, delete)
		diags = append(diags, applyDiags...)
		if diags.HasError() {
			return diags
		}

		// Update mapping table accordingly
		err = m.updateDynamoTable(m.ResourceType, target.OrgId, targetResourceId, delete)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
	}
	return
}

func (m *MrMo) updateDynamoTable(resourceType, targetOrgId, targetResourceId string, delete bool) (err error) {
	if !delete {
		return mockDynamo.UpdateItem(resourceType, m.Id, targetOrgId, targetResourceId)
	}
	return mockDynamo.DeleteItem(m.Id)
}
