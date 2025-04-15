package main

import (
	"context"
	"github.com/charliecon/mr-mo-trial-run/mrmo"
	credentialManager "github.com/charliecon/mr-mo-trial-run/mrmo/org_manager"
	"log"
)

func main() {
	const (
		credsFilePath = "./creds.yml"
		isDelete      = true
	)

	resourceType := "genesyscloud_routing_wrapupcode"
	entityId := "6da92528-0107-4816-963b-cee291c0596c"

	credData, err := credentialManager.ParseCredentialData(credsFilePath)
	if err != nil {
		log.Fatal(err)
	}

	var message = mrmo.Message{
		ResourceType: resourceType,
		EntityId:     entityId,
	}

	err = mrmo.ProcessMessage(context.Background(), message, *credData, isDelete)
	if err != nil {
		log.Fatal(err)
	}
}
