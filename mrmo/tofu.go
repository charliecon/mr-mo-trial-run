package mrmo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"os"
	"os/exec"
)

func runTofu(dir, sourceEntityId, resourcePath string, isDelete bool) (string, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Initialize OpenTofu
	initCmd := exec.Command("tofu", "init")
	initCmd.Dir = dir
	initCmd.Stdout = os.Stdout
	initCmd.Stderr = os.Stderr

	if err := initCmd.Run(); err != nil {
		return "", append(diags, diag.FromErr(err)...)
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
		return "", append(diags, diag.FromErr(err)...)
	}

	if isDelete {
		return "", diags
	}

	targetEntityId, err := extractTargetEntityIdFromOutputs(dir, sourceEntityId)
	if err != nil {
		return "", append(diags, diag.FromErr(err)...)
	}

	return targetEntityId, diags
}

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

	return "", fmt.Errorf("could not find target entity ID in outputs for source entity '%s'", sourceEntityId)
}
