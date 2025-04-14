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

	gcResourceExporter, diags := tfexporter.NewGenesysCloudResourceExporter(ctx, exportResourceConfig, mrMo.ProviderMeta, tfexporter.IncludeResources)
	if diags.HasError() {
		return buildErrorFromDiagnostics(diags)
	}

	exporter := providerRegistrar.GetResourceExporterByResourceType(message.ResourceType)

	m, exportDiags := gcResourceExporter.ExportForMrMo(message.ResourceType, exporter, message.EntityId)
	diags = append(diags, exportDiags...)
	if diags.HasError() {
		return buildErrorFromDiagnostics(diags)
	}

	diags = mrMo.apply(m, false)
	if diags.HasError() {
		return buildErrorFromDiagnostics(diags)
	}

	return nil
}

func (m *MrMo) apply(resourceConfig util.JsonMap, delete bool) (diags diag.Diagnostics) {
	for _, target := range m.OrgManager.Targets {
		// determine if resource file exists for this org
		fm := newFileManager(target.Id, m.Id)

		diags = append(diags, fm.updateTargetTfConfig(resourceConfig, delete)...)
		if diags.HasError() {
			return
		}

		// run targeted apply
		diags = append(diags, runTofu(fm.targetConfigDir)...)
		if diags.HasError() {
			return
		}
	}
	return
}
