package main

import (
	"context"
	"github.com/charliecon/mr-mo-trial-run/mrmo"
	credentialManager "github.com/charliecon/mr-mo-trial-run/mrmo/org_manager"
	"log"
)

func main() {
	const credsFilePath = "./creds.yml"
	var (
		resourceType = "genesyscloud_routing_wrapupcode"
		entityId     = "aba633c3-1ffc-4aa4-84f3-93129b55238a"
		isDelete     = true
	)

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
