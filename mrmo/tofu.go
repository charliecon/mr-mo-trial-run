package mrmo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"os"
	"os/exec"
)

// applyWithOpenTofu executes OpenTofu commands to apply infrastructure changes for a given resource.
//
// Parameters:
//   - dir: The directory path containing the OpenTofu configuration files
//   - sourceEntityId: The identifier of the source entity being processed
//   - resourcePath: The specific resource path to target in the OpenTofu configuration (e.g., "genesyscloud_group.example")
//   - isDelete: Boolean flag indicating if this is a delete operation
//
// Returns:
//   - string: The target entity ID extracted from OpenTofu outputs (empty string if isDelete is true)
//   - diag.Diagnostics: Collection of any diagnostic messages or errors encountered
//
// The function performs the following steps:
//  1. Initializes OpenTofu in the specified directory
//  2. Applies the configuration with appropriate targeting
//  3. For non-delete operations, extracts and returns the target entity ID from outputs
func applyWithOpenTofu(dir, sourceEntityId, resourcePath string, isDelete bool) (_ string, diags diag.Diagnostics) {
	// Initialize OpenTofu
	initCmd := exec.Command("tofu", "init")
	initCmd.Dir = dir
	initCmd.Stdout = os.Stdout
	initCmd.Stderr = os.Stderr

	if err := initCmd.Run(); err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return "", diags
	}

	applyCommand := []string{"apply", "-auto-approve"}
	if !isDelete {
		applyCommand = append(applyCommand, "-target", resourcePath)
	}

	// Apply the configuration
	applyCmd := exec.Command("tofu", applyCommand...)
	applyCmd.Dir = dir
	applyCmd.Stdout = os.Stdout
	applyCmd.Stderr = os.Stderr

	if err := applyCmd.Run(); err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return "", diags
	}

	if isDelete {
		return "", diags
	}

	targetEntityId, err := extractTargetEntityIdFromOutputs(dir, sourceEntityId)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return "", diags
	}

	return targetEntityId, diags
}

// extractTargetEntityIdFromOutputs retrieves and parses the OpenTofu outputs to find the target entity ID
// corresponding to the given source entity.
//
// Parameters:
//   - dir: The directory path containing the OpenTofu configuration and state
//   - sourceEntityId: The identifier of the source entity. This is used as the output label in the target tf configuration, with
//     the value being the ID of the target entity.
//
// Returns:
//   - string: The target entity ID found in the outputs
//   - error: Error if the output command fails, JSON parsing fails, or the target entity ID cannot be found
//
// The function executes 'tofu output -json' and parses the JSON response to find an output block
// that matches the source entity ID. If no matching output is found, returns an error with details about the missing output.
func extractTargetEntityIdFromOutputs(dir, sourceEntityId string) (string, error) {
	var outputBuffer bytes.Buffer

	// Get the output values
	outputCmd := exec.Command("tofu", "output", "-json")
	outputCmd.Dir = dir
	outputCmd.Stdout = &outputBuffer
	outputCmd.Stderr = os.Stderr

	if err := outputCmd.Run(); err != nil {
		return "", err
	}

	// Parse the JSON output
	var outputs map[string]struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(outputBuffer.Bytes(), &outputs); err != nil {
		return "", err
	}

	// Assuming output block is named after the source entity ID
	if output, exists := outputs[buildOutputKey(sourceEntityId)]; exists {
		return output.Value, nil
	}

	return "", fmt.Errorf("could not find target entity ID in outputs for source entity '%s'. Dir: '%s'", sourceEntityId, dir)
}
