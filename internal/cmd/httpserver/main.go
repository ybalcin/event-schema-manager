package main

import (
	"../../../internal/handlers"
	"../../../internal/repositories"
	"../../../pkg/schemaregistry"
	"../../core/services"
)

func main() {
	schemaRegistryClient, err := schemaregistry.NewClient(config.schemaRegistryUrl)
	if err != nil {
		panic(err)
	}

	schemaRegistryAdapter := repositories.NewSchemaRegistryAdapter(schemaRegistryClient)
	schemaService := services.NewSchemaService(schemaRegistryAdapter)
	handlers.NewHttpHandler(schemaService)
}
