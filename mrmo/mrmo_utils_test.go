package mrmo

import (
	"github.com/google/uuid"
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/resource_exporter"
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/tfexporter"
	"testing"
)

func TestUnitExtractGuidsUsingExporterRefAttrs(t *testing.T) {
	entityGuid := uuid.NewString()
	const (
		resourceType  = "resource_foo_bar"
		resourceLabel = "example"
	)
	fullPath := resourceType + "." + resourceLabel

	m := MrMo{
		ResourceType: resourceType,
		Id:           entityGuid,
		Exporter: &resource_exporter.ResourceExporter{
			RefAttrs: map[string]*resource_exporter.RefAttrSettings{
				"base_value_1": {
					RefType: "example_resource",
				},
				"base_value_2": {
					RefType: "example_resource_two",
				},
				"grandparent.parent.grandchild": {
					RefType: "example_resource_three",
				},
				"grandparent.parent.grandchild_2.great_grandchild": {
					RefType: "example_resource_four",
				},
				"not_in_config": {
					RefType: "example_resource_five",
				},
			},
		},
		ResourcePath:  fullPath,
		ResourceLabel: resourceLabel,
	}

	nestedGuid1 := uuid.NewString()
	nestedGuid2 := uuid.NewString()
	nestedGuid3 := uuid.NewString()
	nestedGuid4 := uuid.NewString()
	nestedGuid5 := uuid.NewString()
	guidThatShouldNotBeThere := uuid.NewString()

	config := map[string]any{
		"resource": map[string]tfexporter.ResourceJSONMaps{
			resourceType: {
				resourceLabel: map[string]any{
					"base_value_1":         nestedGuid1,
					"base_value_2":         []any{nestedGuid2},
					"guid_not_in_refattrs": guidThatShouldNotBeThere,
					"grandparent": map[string]any{
						"parent": map[string]any{
							"grandchild": nestedGuid3,
							"grandchild_2": map[string]any{
								"great_grandchild": []string{nestedGuid4, nestedGuid5},
							},
						},
					},
				},
			},
		},
	}

	guids := m.extractGuidsUsingExporterRefAttrs(config)
	validateStringsExistInSlice(t, guids, nestedGuid1, nestedGuid2, nestedGuid3, nestedGuid4, nestedGuid5)
	if stringInStringSlice(guids, guidThatShouldNotBeThere) {
		t.Errorf("Expected %s to not be in slice", guidThatShouldNotBeThere)
	}
}

func validateStringsExistInSlice(t *testing.T, slice []string, s ...string) {
	for _, v := range s {
		if !stringInStringSlice(slice, v) {
			t.Errorf("Expected %s to be in returned slice", v)
		}
	}
}

func stringInStringSlice(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
