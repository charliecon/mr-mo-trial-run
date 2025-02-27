package main

import (
	"context"
	"github.com/charliecon/mr-mo-trial-run/mrmo"
	credentialManager "github.com/charliecon/mr-mo-trial-run/mrmo/credential_manager"
	"log"
)

func main() {
	const credsFilePath = "./creds.yml"
	var (
		resourceType   = "genesyscloud_routing_skill"
		sourceEntityId = "fb9127a3-5ec5-4a1b-abad-79779b48e225"
	)

	credData, err := credentialManager.ParseCredentialData(credsFilePath)
	if err != nil {
		log.Fatal(err)
	}

	var message = mrmo.Message{
		ResourceType: resourceType,
		EntityId:     sourceEntityId,
		Operation:    mrmo.Create,
	}

	err = mrmo.ProcessMessage(context.Background(), message, *credData)
	if err != nil {
		log.Fatal(err)
	}
}
