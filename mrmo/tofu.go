package mrmo

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"os"
	"os/exec"
)

func runTofu(dir, resourcePath string, isDelete bool) (diags diag.Diagnostics) {
	// Initialize OpenTofu
	initCmd := exec.Command("tofu", "init")
	initCmd.Dir = dir
	initCmd.Stdout = os.Stdout
	initCmd.Stderr = os.Stderr

	if err := initCmd.Run(); err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return
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
	}
	return
}
