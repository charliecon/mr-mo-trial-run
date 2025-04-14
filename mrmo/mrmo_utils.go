package mrmo

import (
	"fmt"
	credentialManager "github.com/charliecon/mr-mo-trial-run/mrmo/credential_manager"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mypurecloud/platform-client-sdk-go/v154/platformclientv2"
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/mrmo"
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/provider"
	providerRegistrar "github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/provider_registrar"
	"testing"
)

func newMrMo(resourceType string, orgData credentialManager.CredentialManager) (*MrMo, error) {
	var m MrMo

	m.ResourceType = resourceType

	// initialise ProviderMeta
	providerMeta, err := getProviderConfig(orgData)
	if err != nil {
		return nil, err
	}
	m.ProviderMeta = providerMeta

	mrmo.Activate(providerMeta.ClientConfig)

	// initialise SchemaResource
	allResources, _ := providerRegistrar.GetProviderResources()
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
