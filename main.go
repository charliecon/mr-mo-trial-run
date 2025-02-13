package main

import (
	"log"
	"mr-mo-trial-run/mrmo"
	"time"
)

func main() {
	var (
		resourceType = "genesyscloud_routing_skill"
		data         = map[string]any{
			"name": "Test Routing Skill 1202",
		}
	)

	log.Println("Retrieving Mr Mo instance")
	mrMo, err := mrmo.GetMrMoInstance(resourceType, data)
	if err != nil {
		panic(err)
	}
	log.Println("Successfully retrieved Mr Mo instance")

	log.Printf("Creating %s", resourceType)
	err = mrMo.Create()
	if err != nil {
		panic(err)
	}
	log.Printf("Successfully created %s. ID: %s", resourceType, mrMo.Id)

	log.Println("Sleeping for 2 seconds...")
	time.Sleep(2 * time.Second)

	log.Printf("Deleting %s %s", resourceType, mrMo.Id)
	err = mrMo.Delete()
	if err != nil {
		panic(err)
	}
	log.Printf("Successfully deleted %s %s", resourceType, mrMo.Id)
}
