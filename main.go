package main

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

func main() {
	resourceType := "genesyscloud_routing_skill"
	inputData := map[string]interface{}{
		"name": "Test Skill 1202",
	}

	id, err := performCreate(resourceType, inputData)
	if err != nil {
		panic(err)
	}
	fmt.Println("ID: ", id)
}

func performCreate(resourceType string, data map[string]any) (string, error) {
	providerMeta, err := getProviderConfig()
	if err != nil {
		panic(err)
	}

	allResources, _, _ := provider_registrar.GetAllResources()
	schemaResource, ok := allResources[resourceType]
	if !ok {
		return "", fmt.Errorf("resource not found %s", resourceType)
	}

	resourceDataObject := createResourceDataObject(schemaResource.Schema, data)

	createFunc := schemaResource.CreateContext

	diagErr := createFunc(context.Background(), resourceDataObject, providerMeta)
	if diagErr != nil {
		return "", fmt.Errorf("%v", diagErr)
	}
	return resourceDataObject.Id(), nil
}

func getProviderConfig() (*provider.ProviderMeta, error) {
	var (
		clientId     = os.Getenv("GENESYSCLOUD_OAUTHCLIENT_ID")
		clientSecret = os.Getenv("GENESYSCLOUD_OAUTHCLIENT_SECRET")
	)

	if clientId == "" || clientSecret == "" {
		return nil, fmt.Errorf("GENESYSCLOUD_OAUTHCLIENT_ID and GENESYSCLOUD_OAUTHCLIENT_SECRET must be set")
	}

	config := platformclientv2.GetDefaultConfiguration()
	err := config.AuthorizeClientCredentials(clientId, clientSecret)
	if err != nil {
		return nil, err
	}

	return &provider.ProviderMeta{
		ClientConfig: config,
	}, nil
}

func createResourceDataObject(routingSkillSchema map[string]*schema.Schema, data map[string]any) *schema.ResourceData {
	var t testing.T
	s := schema.TestResourceDataRaw(&t, routingSkillSchema, data)
	return s
}
