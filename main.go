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
		//routingSkillResourceType   = "genesyscloud_routing_skill"
		//sourceRoutingSKillEntityId = "fb9127a3-5ec5-4a1b-abad-79779b48e225"

		groupResourceType = "genesyscloud_group"
		groupEntityId     = "d6c70405-1351-49bf-a9fe-ec4ba2363ad2"
	)

	credData, err := credentialManager.ParseCredentialData(credsFilePath)
	if err != nil {
		log.Fatal(err)
	}

	var message = mrmo.Message{
		ResourceType: groupResourceType,
		EntityId:     groupEntityId,
		Operation:    mrmo.Create,
	}

	err = mrmo.ProcessMessage(context.Background(), message, *credData)
	if err != nil {
		log.Fatal(err)
	}
}
