## Mr Mo (Prototype)

Mr. Mo (Multi-Region Multi-Org) is a service that enables seamless resource replication across multiple Genesys Cloud organizations. 
It automatically exports, transforms, and applies resource configurations from a source organization to multiple target organizations while maintaining proper GUID mappings and credentials. 
This tool is particularly useful for enterprises managing identical infrastructure across different organizational boundaries, ensuring consistent resource deployment and configuration management at scale.

### Usage

1. Define a credentials file with the same structure as `creds_example.yml`.

2. Build a message object and invoke the `ProcessMessage` function.

Example:

```go
package main

import (
	"context"
	"github.com/charliecon/mr-mo-trial-run/mrmo"
	"log"
)

func main() {
	credsFilePath := "./path/to/your/creds_file.yml"
	
	message := mrmo.Message{
		ResourceType: "genesyscloud_outbound_attempt_limit", // terraform resource type
		EntityId:     "25c1c87e-b7b8-45a8-aab2-fe42295f16a3", // attempt limit ID in the source org
		IsDelete:     false, // was the attempt limit created or deleted in the source org
	}

	diags := mrmo.ProcessMessage(context.Background(), message, credsFilePath)
	if diags.HasError() {
		log.Fatal(diags)
	}
	log.Println(diags)
}
```
