package mrmo

import (
	"context"
	orgManager "github.com/charliecon/mr-mo-trial-run/mrmo/org_manager"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	providerRegistrar "github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/provider_registrar"
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/tfexporter"
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
}

type Message struct {
	ResourceType string
	EntityId     string
}

func ProcessMessage(ctx context.Context, message Message, om orgManager.OrgManager, delete bool) error {
	var diags = make(diag.Diagnostics, 0)
	defer func() {
		printDiagnosticWarnings(diags)
	}()

	mrMo, err := newMrMo(message.ResourceType, om, message.EntityId)
	if err != nil {
		log.Println("Failed to initialise mr mo")
		return err
	}

	if delete {
		diags = append(diags, mrMo.apply(nil, true)...)
		if diags.HasError() {
			return buildErrorFromDiagnostics(diags)
		}
		return nil
	}

	exportResourceConfig := createExportResourceData(tfexporter.ResourceTfExport().Schema, tfexporter.ResourceType)

	gcResourceExporter, newExporterDiags := tfexporter.NewGenesysCloudResourceExporter(ctx, exportResourceConfig, mrMo.ProviderMeta, tfexporter.IncludeResources)
	diags = append(diags, newExporterDiags...)
	if diags.HasError() {
		return buildErrorFromDiagnostics(diags)
	}

	exporter := providerRegistrar.GetResourceExporterByResourceType(message.ResourceType)

	m, exportDiags := gcResourceExporter.ExportForMrMo(message.ResourceType, exporter, message.EntityId)
	diags = append(diags, exportDiags...)
	if diags.HasError() {
		return buildErrorFromDiagnostics(diags)
	}

	resourcePath, err := parseResourcePathFromConfig(m, message.ResourceType)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return buildErrorFromDiagnostics(diags)
	}
	mrMo.ResourcePath = resourcePath

	diags = append(diags, mrMo.apply(m, false)...)
	if diags.HasError() {
		return buildErrorFromDiagnostics(diags)
	}

	return nil
}

func (m *MrMo) apply(resourceConfig util.JsonMap, delete bool) (diags diag.Diagnostics) {
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
			break
		}

		// determine if resource file exists for this org
		fm := newFileManager(target.Id, m.Id)

		diags = append(diags, fm.updateTargetTfConfig(resourceConfig, delete)...)
		if diags.HasError() {
			break
		}

		// run targeted apply
		diags = append(diags, runTofu(fm.targetConfigDir, m.ResourcePath, delete)...)
		if diags.HasError() {
			break
		}
	}
	return
}
