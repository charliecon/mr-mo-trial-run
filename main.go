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
	knowledgeBaseResourceType := "genesyscloud_knowledge_knowledgebase"
	knowledgeBaseId := "314ba3df-1bd4-4681-b250-2207a6d97bc9"

	//knowledgeDocumentResourceType := "genesyscloud_knowledge_document"
	//documentId := "a9ae73da-eca6-459a-814e-8242c55edc9e," + knowledgeBaseId

	credData, err := orgManager.ParseCredentialData(credsFilePath)
	if err != nil {
		log.Fatal(err)
	}

	var message = mrmo.Message{
		ResourceType: knowledgeBaseResourceType,
		EntityId:     knowledgeBaseId,
		IsDelete:     true,
	}

	diags := mrmo.ProcessMessage(ctx, message, *credData)
	if diags.HasError() {
		log.Fatal(diags)
	}
}
