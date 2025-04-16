package main

import (
	"context"
	"github.com/charliecon/mr-mo-trial-run/mrmo"
	orgManager "github.com/charliecon/mr-mo-trial-run/mrmo/org_manager"
	"log"
)

func main() {
	const credsFilePath = "./creds.yml"

	ctx := context.Background()
	resourceType := "genesyscloud_knowledge_knowledgebase"
	entityId := "314ba3df-1bd4-4681-b250-2207a6d97bc9"

	credData, err := orgManager.ParseCredentialData(credsFilePath)
	if err != nil {
		log.Fatal(err)
	}

	var message = mrmo.Message{
		ResourceType: resourceType,
		EntityId:     entityId,
		IsDelete:     true,
	}

	diags := mrmo.ProcessMessage(ctx, message, *credData)
	if diags.HasError() {
		log.Fatal(diags)
	}
}
