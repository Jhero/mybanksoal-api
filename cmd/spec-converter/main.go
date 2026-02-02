package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"gopkg.in/yaml.v3"
)

func main() {
	// Read Swagger 2.0 file
	data, err := os.ReadFile("docs/swagger.json")
	if err != nil {
		log.Fatalf("Failed to read swagger.json: %v", err)
	}

	// Parse Swagger 2.0
	var doc2 openapi2.T
	if err := json.Unmarshal(data, &doc2); err != nil {
		log.Fatalf("Failed to parse swagger.json: %v", err)
	}

	// Convert to OpenAPI 3.0
	doc3, err := openapi2conv.ToV3(&doc2)
	if err != nil {
		log.Fatalf("Failed to convert to OpenAPI 3.0: %v", err)
	}

	// Write OpenAPI 3.0 file
	out, err := yaml.Marshal(doc3)
	if err != nil {
		log.Fatalf("Failed to marshal OpenAPI 3.0: %v", err)
	}

	if err := os.WriteFile("api.yml", out, 0644); err != nil {
		log.Fatalf("Failed to write api.yml: %v", err)
	}

	log.Println("Successfully converted Swagger 2.0 to OpenAPI 3.0 (api.yml)")
}
