package mrmo

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/tfexporter"
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/util"
	"log"
	"os"
	"path/filepath"
)

type FileManager struct {
	targetOrgId      string
	sourceEntityId   string
	targetConfigDir  string
	targetConfigFile string
	exists           bool
}

func newFileManager(targetOrgId, sourceEntityId string) *FileManager {
	fm := FileManager{
		targetOrgId:      targetOrgId,
		sourceEntityId:   sourceEntityId,
		targetConfigDir:  filepath.Join("mock-s3", "organizations", targetOrgId, "config"),
		targetConfigFile: filepath.Join("mock-s3", "organizations", targetOrgId, "config", sourceEntityId+".tf.json"),
	}
	fm.exists = fileExists(fm.targetConfigFile)
	return &fm
}

func (f *FileManager) updateTargetTfConfig(resourceConfig util.JsonMap, delete bool) (diags diag.Diagnostics) {
	if delete {
		diags = append(diags, f.deleteResourceFile()...)
		return
	}
	return append(diags, tfexporter.WriteConfigForMrMo(resourceConfig, f.targetConfigFile)...)
}

func (f *FileManager) deleteResourceFile() (diags diag.Diagnostics) {
	if !f.exists {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "File did not exist before delete was invoked.",
			Detail:   fmt.Sprintf("We are not attempting to delete file '%s' because it already does not exist.", f.targetConfigFile),
		})
		return
	}

	log.Printf("Deleting file '%s'", f.targetConfigFile)
	err := os.Remove(f.targetConfigFile)
	if err == nil {
		return
	}

	if os.IsNotExist(err) {
		log.Printf("Successfully deleted file '%s'", f.targetConfigFile)
		return
	}
	return append(diags, diag.Errorf("failed to delete file '%s'. Error: %s", f.targetConfigFile, err.Error())...)
}

func fileExists(filePath string) bool {
	var err error
	if _, err = os.Stat(filePath); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	}
	log.Printf("Failed to verify if file '%s' exists. Returning false. Error: '%s'", filePath, err.Error())
	return false
}
