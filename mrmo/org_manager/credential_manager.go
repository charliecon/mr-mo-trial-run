package org_manager

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

type OrgManager struct {
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

func ParseCredentialData(credsFilePath string) (*OrgManager, error) {
	filename, err := filepath.Abs(credsFilePath)
	if err != nil {
		return nil, err
	}

	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var secretManager OrgManager
	err = yaml.Unmarshal(yamlFile, &secretManager)
	if err != nil {
		return nil, err
	}

	return &secretManager, nil
}

const (
	clientIdEnvVar     = "GENESYSCLOUD_OAUTHCLIENT_ID"
	clientSecretEnvVar = "GENESYSCLOUD_OAUTHCLIENT_SECRET"
	regionEnvVar       = "GENESYSCLOUD_REGION"
)

func (o *OrgData) SetTargetOrgCredentials() error {
	return SetClientCredEnvVars(o.ClientId, o.ClientSecret, o.Region)
}

func SetClientCredEnvVars(id, secret, region string) (err error) {
	setEnvVarFunc := func(envVar, value string) error {
		if setErr := os.Setenv(envVar, value); err != nil {
			return fmt.Errorf("failed to set env var '%s' to value '%s'. Error: %w", envVar, value, setErr)
		}
		return nil
	}

	if err = setEnvVarFunc(clientIdEnvVar, id); err != nil {
		return err
	}
	if err = setEnvVarFunc(clientSecretEnvVar, secret); err != nil {
		return err
	}
	return setEnvVarFunc(regionEnvVar, region)
}

func GetClientCredsEnvVars() (id, secret, region string) {
	return os.Getenv(clientIdEnvVar), os.Getenv(clientSecretEnvVar), os.Getenv(regionEnvVar)
}
