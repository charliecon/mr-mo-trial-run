package main

import (
	"fmt"
	"mr-mo-trial-run/mrmo"
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

	err = mrMo.Create()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Successfully created %s. ID: %s", resourceType, mrMo.Id)
}
