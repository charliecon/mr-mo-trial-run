package main

import (
	"context"
	"fmt"
	"github.com/charliecon/mr-mo-trial-run/mrmo"
	orgManager "github.com/charliecon/mr-mo-trial-run/mrmo/org_manager"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
)

const credsFilePath = "./creds.yml"
const attemptLimitResourceType = "genesyscloud_outbound_attempt_limit"
const cltResourceType = "genesyscloud_outbound_contact_list_template"

func main() {
	var (
		ctx   = context.Background()
		diags diag.Diagnostics
	)

	defer func() {
		printDiagnosticWarnings(diags)
	}()

	credData, err := orgManager.ParseCredentialData(credsFilePath)
	if err != nil {
		log.Fatal(err)
	}

	delete := true

	attemptLimitId := "e6c22ee5-6dff-4da9-8344-e985a1a269e4"
	cltId := "87725a8f-b66d-45ef-9631-e816747ad5b7"

	// process attempt limit
	message := mrmo.Message{
		ResourceType: attemptLimitResourceType,
		EntityId:     attemptLimitId,
		IsDelete:     delete,
	}

	diags = append(diags, mrmo.ProcessMessage(ctx, message, *credData)...)
	if diags.HasError() {
		log.Fatal(diags)
	}

	// process contact list template
	message = mrmo.Message{
		ResourceType: cltResourceType,
		EntityId:     cltId,
		IsDelete:     delete,
	}

	diags = append(diags, mrmo.ProcessMessage(ctx, message, *credData)...)
	if diags.HasError() {
		log.Fatal(diags)
	}
}

// printDiagnosticWarnings will print any diagnostics warnings, if any exist
func printDiagnosticWarnings(diags diag.Diagnostics) {
	if len(diags) == 0 || diags.HasError() {
		return
	}
	log.Println("Diagnostic warnings: ")
	for _, d := range diags {
		fmt.Println(d)
	}
}
