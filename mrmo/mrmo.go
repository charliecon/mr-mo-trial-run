package mrmo

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mypurecloud/platform-client-sdk-go/v150/platformclientv2"
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/provider"
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/provider_registrar"
	"os"
	"testing"
)

var mrMoInstance *MrMo

type MrMo struct {
	ResourceType   string
	Id             string
	Data           map[string]any
	ResourceData   *schema.ResourceData
	SchemaResource *schema.Resource
	ProviderMeta   any
}

func (m *MrMo) Create() error {
	diagErr := m.SchemaResource.CreateContext(context.Background(), m.ResourceData, m.ProviderMeta)
	if diagErr != nil {
		return fmt.Errorf("%v", diagErr)
	}
	m.Id = m.ResourceData.Id()
	return nil
}

func (m *MrMo) Delete() error {
	diagErr := m.SchemaResource.DeleteContext(context.Background(), m.ResourceData, m.ProviderMeta)
	if diagErr != nil {
		return fmt.Errorf("%v", diagErr)
	}
	return nil
}

func GetMrMoInstance(resourceType string, data map[string]any) (_ *MrMo, err error) {
	if mrMoInstance == nil {
		mrMoInstance, err = newMrMo(resourceType, data)
	}
	return mrMoInstance, err
}

func newMrMo(resourceType string, data map[string]any) (*MrMo, error) {
	var m MrMo

	m.ResourceType = resourceType
	m.Data = data

	// initialise ProviderMeta
	providerMeta, err := getProviderConfig()
	if err != nil {
		return nil, err
	}
	m.ProviderMeta = providerMeta

	// initialise SchemaResource
	allResources, _, _ := provider_registrar.GetAllResources()
	schemaResource, ok := allResources[resourceType]
	if !ok {
		return nil, fmt.Errorf("resource not found %s", resourceType)
	}
	m.SchemaResource = schemaResource

	// initialise SchemaResource
	resourceDataObject := createResourceDataObject(schemaResource.Schema, data)
	m.ResourceData = resourceDataObject

	return &m, nil
}

func createResourceDataObject(routingSkillSchema map[string]*schema.Schema, data map[string]any) *schema.ResourceData {
	var t testing.T
	return schema.TestResourceDataRaw(&t, routingSkillSchema, data)
}

func getProviderConfig() (_ *provider.ProviderMeta, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("getProviderConfig: %w", err)
		}
	}()

	var (
		clientId     = os.Getenv("GENESYSCLOUD_OAUTHCLIENT_ID")
		clientSecret = os.Getenv("GENESYSCLOUD_OAUTHCLIENT_SECRET")
	)

	if clientId == "" || clientSecret == "" {
		return nil, fmt.Errorf("GENESYSCLOUD_OAUTHCLIENT_ID and GENESYSCLOUD_OAUTHCLIENT_SECRET must be set")
	}

	config := platformclientv2.GetDefaultConfiguration()
	err = config.AuthorizeClientCredentials(clientId, clientSecret)
	if err != nil {
		return nil, err
	}

	return &provider.ProviderMeta{
		ClientConfig: config,
	}, nil
}
