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

	credData, err := orgManager.ParseCredentialData(credsFilePath)
	if err != nil {
		log.Fatal(err)
	}

	processClt := false
	delete := true

	attemptLimitResourceType := "genesyscloud_outbound_attempt_limit"
	attemptLimitId := "e6c22ee5-6dff-4da9-8344-e985a1a269e4"

	cltResourceType := "genesyscloud_outbound_contact_list_template"
	cltId := "87725a8f-b66d-45ef-9631-e816747ad5b7"

	var message mrmo.Message

	if processClt {
		message = mrmo.Message{
			ResourceType: cltResourceType,
			EntityId:     cltId,
			IsDelete:     delete,
		}
	} else {
		message = mrmo.Message{
			ResourceType: attemptLimitResourceType,
			EntityId:     attemptLimitId,
			IsDelete:     delete,
		}
	}

	diags := mrmo.ProcessMessage(ctx, message, *credData)
	if diags.HasError() {
		log.Fatal(diags)
	}
}
