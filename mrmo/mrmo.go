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
	ResourcePath   string // determined after export
	ProviderMeta   any
	OrgManager     *orgManager.OrgManager
	Exporter       *resourceExporter.ResourceExporter
}

type Message struct {
	ResourceType string
	EntityId     string
	IsDelete     bool
}

func ProcessMessage(ctx context.Context, message Message, om orgManager.OrgManager) (diags diag.Diagnostics) {
	mrMo, err := newMrMo(message.ResourceType, om, message.EntityId)
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

	resourcePath, err := parseResourcePathFromConfig(resourceConfig, message.ResourceType)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	mrMo.ResourcePath = resourcePath

	resourceConfig = appendOutputBlockToConfig(resourceConfig, resourcePath, message.EntityId)

	diags = append(diags, mrMo.applyResourceConfigToTargetOrgs(resourceConfig, false)...)
	return diags
}

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
		err := target.SetTargetOrgCredentials()
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}

		resourceConfigAfterResolvingGuids, err := m.resolveResourceConfigDependencies(resourceConfig, target)
		if err != nil {
			return diag.FromErr(err)
		}

		fm := newFileManager(target.Id, m.Id)

		diags = append(diags, fm.updateTargetTfConfig(resourceConfigAfterResolvingGuids, delete)...)
		if diags.HasError() {
			return diags
		}

		// run targeted apply
		targetResourceId, applyDiags := applyWithOpenTofu(fm.targetConfigDir, m.Id, m.ResourcePath, delete)
		diags = append(diags, applyDiags...)
		if diags.HasError() {
			return diags
		}

		err = m.updateDynamoTable(target.Id, targetResourceId, delete)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
	}
	return
}

func (m *MrMo) updateDynamoTable(targetOrgId, targetResourceId string, delete bool) (err error) {
	if !delete {
		return mockDynamo.UpdateItem(m.Id, targetOrgId, targetResourceId)
	}
	return mockDynamo.DeleteItem(m.Id)
}
