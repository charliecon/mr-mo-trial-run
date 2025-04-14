package mrmo

import (
	"context"
	"fmt"
	credentialManager "github.com/charliecon/mr-mo-trial-run/mrmo/credential_manager"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	providerRegistrar "github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/provider_registrar"
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/tfexporter"
	"log"
	"testing"
)

var mrMoInstance *MrMo

type MrMo struct {
	ResourceType   string
	Id             string
	ResourceData   *schema.ResourceData
	SchemaResource *schema.Resource
	ProviderMeta   any
}

type Operation string

const (
	Create Operation = "Create"
	Update Operation = "Update"
	Delete Operation = "Delete"
)

type Message struct {
	ResourceType string
	EntityId     string
	Operation    Operation
}

func ProcessMessage(ctx context.Context, message Message, orgData credentialManager.CredentialManager) error {
	fmt.Println(orgData.Source.ClientId, orgData.Source.ClientSecret)
	mrMo, err := newMrMo(message.ResourceType, orgData)
	if err != nil {
		log.Println("Failed to initialise mr mo")
		return err
	}

	rd := createExportResourceData(tfexporter.ResourceTfExport().Schema, tfexporter.ResourceType)

	gcResourceExporter, diags := tfexporter.NewGenesysCloudResourceExporter(ctx, rd, mrMo.ProviderMeta, tfexporter.IncludeResources)
	if diags.HasError() {
		return fmt.Errorf("%v", diags)
	}

	exporter := providerRegistrar.GetResourceExporterByResourceType(message.ResourceType)

	m, diags := gcResourceExporter.ExportForMrMo(message.ResourceType, exporter, message.EntityId)
	if diags.HasError() {
		return fmt.Errorf("%v", diags)
	}

	diags = append(diags, tfexporter.WriteConfigForMrMo(m, "./please.tf.json")...)
	if diags.HasError() {
		return fmt.Errorf("%v", diags)
	}
	return nil
}

func createExportResourceData(s map[string]*schema.Schema, resType string) *schema.ResourceData {
	config := map[string]any{
		"include_state_file":       true,
		"export_format":            "json",
		"include_filter_resources": []any{resType},
	}

	var t testing.T
	return schema.TestResourceDataRaw(&t, s, config)
}

func (m *MrMo) Create(ctx context.Context) error {
	diagErr := m.SchemaResource.CreateContext(ctx, m.ResourceData, m.ProviderMeta)
	if diagErr != nil {
		return fmt.Errorf("%v", diagErr)
	}
	m.Id = m.ResourceData.Id()
	return nil
}

func (m *MrMo) Read(ctx context.Context) error {
	diagErr := m.SchemaResource.ReadContext(ctx, m.ResourceData, m.ProviderMeta)
	if diagErr != nil {
		return fmt.Errorf("%v", diagErr)
	}
	m.Id = m.ResourceData.Id()
	return nil
}

func (m *MrMo) Delete(ctx context.Context) error {
	diagErr := m.SchemaResource.DeleteContext(ctx, m.ResourceData, m.ProviderMeta)
	if diagErr != nil {
		return fmt.Errorf("%v", diagErr)
	}
	return nil
}

/*
	instanceState := &terraform.InstanceState{
		ID: message.EntityId,
	}

	state, diagErr := mrMo.SchemaResource.RefreshWithoutUpgrade(ctx, instanceState, mrMo.ProviderMeta)
	if diagErr.HasError() {
		return fmt.Errorf("%v", diagErr)
	}

	resourceInfo := resource_exporter.ResourceInfo{
		State:         state,
		BlockLabel:    "hello",
		OriginalLabel: "hello-og",
		Type:          message.ResourceType,
		BlockType:     "",
		CtyType:       mrMo.SchemaResource.CoreConfigSchema().ImpliedType(),
	}

	fmt.Println("Here is the state:")
	fmt.Println(state.String())
	fmt.Println("here are the attributes:")
	fmt.Println(state.Attributes)
*/
