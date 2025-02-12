package main

import (
	"log"
	"mr-mo-trial-run/mrmo"
	"time"
)

func main() {
	var (
		mrMo mrmo.MrMo

		resourceType = "genesyscloud_routing_skill"
		data         = map[string]any{
			"name": "Test Routing Skill 1202",
		}
	)

	err := mrMo.InitMrMo(resourceType, data)
	if err != nil {
		panic(err)
	}

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
