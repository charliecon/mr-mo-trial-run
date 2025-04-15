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
	entityId := "ee5052a5-ebab-4d01-93f1-663602d64a5f"

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
