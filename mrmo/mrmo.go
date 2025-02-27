package mrmo

import (
	"context"
	"fmt"
	credentialManager "github.com/charliecon/mr-mo-trial-run/mrmo/credential_manager"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mypurecloud/platform-client-sdk-go/v150/platformclientv2"
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/provider"
	providerRegistrar "github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/provider_registrar"
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
	mrMo.ResourceData.SetId(message.EntityId)

	// read from source org
	err = mrMo.Read(ctx)
	if err != nil {
		log.Println("Failed to read from source org")
		return err
	}

	fmt.Println(mrMo.ResourceData.Get("name").(string))

	// perform operation in target orgs

	return nil
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

func newMrMo(resourceType string, orgData credentialManager.CredentialManager) (*MrMo, error) {
	var m MrMo

	m.ResourceType = resourceType

	// initialise ProviderMeta
	providerMeta, err := getProviderConfig(orgData)
	if err != nil {
		return nil, err
	}
	m.ProviderMeta = providerMeta

	// initialise SchemaResource
	allResources, _, _ := providerRegistrar.GetAllResources()
	schemaResource, ok := allResources[resourceType]
	if !ok {
		return nil, fmt.Errorf("resource not found %s", resourceType)
	}
	m.SchemaResource = schemaResource

	// initialise SchemaResource
	resourceDataObject := createResourceDataObject(schemaResource.Schema, make(map[string]any))
	m.ResourceData = resourceDataObject

	return &m, nil
}

func createResourceDataObject(resourceSchema map[string]*schema.Schema, data map[string]any) *schema.ResourceData {
	var t testing.T
	return schema.TestResourceDataRaw(&t, resourceSchema, data)
}

func getProviderConfig(orgData credentialManager.CredentialManager) (_ *provider.ProviderMeta, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("getProviderConfig: %w", err)
		}
	}()

	config := platformclientv2.GetDefaultConfiguration()
	config.BasePath = provider.GetRegionBasePath(orgData.Source.Region)

	err = config.AuthorizeClientCredentials(orgData.Source.ClientId, orgData.Source.ClientSecret)
	if err != nil {
		return nil, err
	}

	return &provider.ProviderMeta{
		ClientConfig: config,
	}, nil
}
