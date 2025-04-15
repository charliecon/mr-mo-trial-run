package main

import (
	"context"
	"github.com/charliecon/mr-mo-trial-run/mrmo"
	credentialManager "github.com/charliecon/mr-mo-trial-run/mrmo/org_manager"
	"log"
)

func main() {
	const credsFilePath = "./creds.yml"

	resourceType := "genesyscloud_routing_wrapupcode"
	entityId := "b039fe91-33e0-4f63-91fd-c1e164f21abe"

	credData, err := credentialManager.ParseCredentialData(credsFilePath)
	if err != nil {
		log.Fatal(err)
	}

	var message = mrmo.Message{
		ResourceType: resourceType,
		EntityId:     entityId,
		IsDelete:     true,
	}

	err = mrmo.ProcessMessage(context.Background(), message, *credData)
	if err != nil {
		log.Fatal(err)
	}
}
