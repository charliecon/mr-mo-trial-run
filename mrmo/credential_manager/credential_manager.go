package credential_manager

import (
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

type CredentialManager struct {
	Source  OrgData   `yaml:"source"`
	Targets []OrgData `yaml:"targets"`
}

type OrgData struct {
	Id           string `yaml:"orgId"`
	Name         string `yaml:"orgName"`
	ClientId     string `yaml:"clientId"`
	ClientSecret string `yaml:"clientSecret"`
	Region       string `yaml:"region"`
}

func ParseCredentialData(credsFilePath string) (*CredentialManager, error) {
	filename, err := filepath.Abs(credsFilePath)
	if err != nil {
		return nil, err
	}

	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var secretManager CredentialManager
	err = yaml.Unmarshal(yamlFile, &secretManager)
	if err != nil {
		return nil, err
	}

	return &secretManager, nil
}
