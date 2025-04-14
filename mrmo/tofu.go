package mrmo

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"os"
	"os/exec"
)

func runTofu(dir string) (diags diag.Diagnostics) {
	// do not do anything yet
	if true {
		return
	}

	// Initialize OpenTofu
	initCmd := exec.Command("tofu", "init")
	initCmd.Dir = dir
	initCmd.Stdout = os.Stdout
	initCmd.Stderr = os.Stderr

	if err := initCmd.Run(); err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return
	}

	// Apply the configuration
	applyCmd := exec.Command("tofu", "apply", "-auto-approve")
	applyCmd.Dir = dir
	applyCmd.Stdout = os.Stdout
	applyCmd.Stderr = os.Stderr

	if err := applyCmd.Run(); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	return
}
