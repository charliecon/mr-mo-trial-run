package mrmo

import (
	"context"
	"encoding/json"
	"fmt"
	mock_dynamo "github.com/charliecon/mr-mo-trial-run/mock-dynamo"
	credentialManager "github.com/charliecon/mr-mo-trial-run/mrmo/org_manager"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mypurecloud/platform-client-sdk-go/v154/platformclientv2"
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/mrmo"
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/provider"
	providerRegistrar "github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/provider_registrar"
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/tfexporter"
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/util"
	"log"
	"regexp"
	"strings"
	"testing"
)

func newMrMo(resourceType string, orgData credentialManager.OrgManager, sourceEntityId string) (*MrMo, error) {
	var m MrMo

	m.ResourceType = resourceType
	m.OrgManager = &orgData
	m.Id = sourceEntityId

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

func getProviderConfig(orgData credentialManager.OrgManager) (_ *provider.ProviderMeta, err error) {
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

// createExportResourceData generates the export resource config that the genesyscloud tf exporter will use
func createExportResourceData(s map[string]*schema.Schema, resType string) *schema.ResourceData {
	config := map[string]any{
		"include_state_file":       true,
		"export_format":            "json",
		"include_filter_resources": []any{resType},
	}

	var t testing.T
	return schema.TestResourceDataRaw(&t, s, config)
}

// printDiagnosticWarnings will print any diagnostics warnings, if any exist
func printDiagnosticWarnings(diags diag.Diagnostics) {
	if len(diags) == 0 || diags.HasError() {
		return
	}
	log.Println("Diagnostic warnings: ")
	for _, d := range diags {
		fmt.Println(d)
	}
}

// parseResourcePathFromConfig will parse the full resource path from the exported resource config
func parseResourcePathFromConfig(resourceConfig util.JsonMap, resourceType string) (_ string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to parse resource path from exported resource config: %w", err)
		}
	}()

	r, ok := resourceConfig["resource"].(any)
	if !ok {
		return "", fmt.Errorf("no resource block found in resource config")
	}

	resourceMap, ok := r.(map[string]tfexporter.ResourceJSONMaps)
	if !ok {
		return "", fmt.Errorf("failed to cast resource map. Expected map[string]any, got %T", r)
	}

	configMap, ok := resourceMap[resourceType]
	if !ok || configMap == nil || len(configMap) == 0 {
		return "", fmt.Errorf("failed to parse config for resource type '%s'", resourceType)
	}

	for resourceLabel, _ := range configMap {
		return resourceType + "." + resourceLabel, nil
	}

	return "", fmt.Errorf("no resource label found for resource type '%s'", resourceType)
}

func (m *MrMo) exportConfig(ctx context.Context, resourceId, resourceType string) (_ util.JsonMap, diags diag.Diagnostics) {
	exportResourceConfig := createExportResourceData(tfexporter.ResourceTfExport().Schema, tfexporter.ResourceType)

	gcResourceExporter, newExporterDiags := tfexporter.NewGenesysCloudResourceExporter(ctx, exportResourceConfig, m.ProviderMeta, tfexporter.IncludeResources)
	diags = append(diags, newExporterDiags...)
	if diags.HasError() {
		return nil, diags
	}

	exporter := providerRegistrar.GetResourceExporterByResourceType(resourceType)
	m.Exporter = exporter

	config, exportDiags := gcResourceExporter.ExportForMrMo(resourceType, exporter, resourceId)
	diags = append(diags, exportDiags...)

	return config, diags
}

// resolveResourceConfigDependencies will find GUIDS inside the resource config and
// try to resolve them to GUIDs in the target org
func (m *MrMo) resolveResourceConfigDependencies(resourceConfig util.JsonMap, target credentialManager.OrgData) (_ util.JsonMap, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("resolveResourceConfigDependencies: %w", err)
		}
	}()

	guidReferencesInConfig, err := extractUUIDs(resourceConfig)
	if err != nil {
		return nil, err
	}

	for _, guid := range guidReferencesInConfig {
		// if we matched the source entity ID - continue. It is there as the name of an output variables
		if guid == m.Id {
			continue
		}
		// search for guid.target.Id value
		item, err := mock_dynamo.GetItem(guid)
		if err != nil {
			log.Printf("Failed to read guid '%s' from dynamo. Error: %s", guid, err.Error())
			continue
		}

		targetGuid := item[target.Id]

		// replace guid with that value
		resourceConfig, err = replaceGUID(resourceConfig, guid, targetGuid)
		if err != nil {
			return nil, err
		}
	}

	return resourceConfig, err
}

// extractUUIDs converts the input map to a string and finds all UUIDs
func extractUUIDs(data util.JsonMap) (_ []string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("extractUUIDs: %w", err)
		}
	}()

	// Convert the data structure to JSON string
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Convert to string
	jsonStr := string(jsonBytes)

	// Regular expression for UUID pattern (excluding the output variables which are prefixed with "{outputPrefix}")
	uuidRegex := regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)

	return uuidRegex.FindAllString(jsonStr, -1), nil
}

// replaceGUID takes a JsonMap, finds all instances of oldGUID within the data structure
// and replaces them with newGUID. The function converts the map to a JSON string to perform
// the replacement, then converts it back to a JsonMap. It returns the modified JsonMap and
// any error encountered during the process.
func replaceGUID(data util.JsonMap, oldGUID string, newGUID string) (_ util.JsonMap, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("replaceGUID: %w", err)
		}
	}()

	// Convert the data structure to JSON string
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling data: %w", err)
	}

	// Convert to string
	jsonStr := string(jsonBytes)

	// Replace the GUID
	updatedStr := strings.Replace(jsonStr, oldGUID, newGUID, -1)

	// Convert back to map
	var result util.JsonMap
	if err := json.Unmarshal([]byte(updatedStr), &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling data: %w", err)
	}

	return result, err
}

func appendOutputBlockToConfig(config util.JsonMap, resourcePath, sourceEntityId string) util.JsonMap {
	config["output"] = map[string]map[string]string{
		buildOutputKey(sourceEntityId): {
			"value": fmt.Sprintf("${%s.id}", resourcePath),
		},
	}
	return config
}

const outputPrefix = "mrmo_"

// buildOutputKey will build the key of the output tf block. They must not start with a number, so the GUID alone will not do.
func buildOutputKey(sourceEntityId string) string {
	return outputPrefix + sourceEntityId
}
