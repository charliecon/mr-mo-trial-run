package main

import (
	"context"
	"fmt"
	"github.com/charliecon/mr-mo-trial-run/mrmo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
)

const (
	credsFilePath            = "./creds.yml"
	attemptLimitResourceType = "genesyscloud_outbound_attempt_limit"
	cltResourceType          = "genesyscloud_outbound_contact_list_template"
)

func main() {
	var (
		ctx   = context.Background()
		diags diag.Diagnostics
	)

	defer func() {
		printDiagnosticWarnings(diags)
	}()

	delete := true

	// process attempt limit
	attemptLimitId := "dfada4b5-2293-4ffa-9af7-b2e53fadbdc4"
	diags = append(diags, processResource(ctx, attemptLimitResourceType, attemptLimitId, delete)...)
	if diags.HasError() {
		log.Fatal(diags)
	}

	// process contact list template
	//cltId := "4ff035cb-8585-4880-99f7-5e414db62946"
	//diags = append(diags, processResource(ctx, cltResourceType, cltId, delete)...)
	//if diags.HasError() {
	//	log.Fatal(diags)
	//}
}

func processResource(ctx context.Context, resourceType, entityId string, delete bool) diag.Diagnostics {
	message := mrmo.Message{
		ResourceType: resourceType,
		EntityId:     entityId,
		IsDelete:     delete,
	}
	return mrmo.ProcessMessage(ctx, message, credsFilePath)
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
