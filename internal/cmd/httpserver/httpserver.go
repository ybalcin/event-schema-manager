package httpserver

import (
	"../../../internal/handlers"
	"../../../internal/repositories"
	"../../../pkg/schemaregistry"
	"../../core/services"
)

func init() {

	schemaRegistryClient, err := schemaregistry.NewClient(config.schemaRegistryUrl)
	if err != nil {
		panic(err)
	}

	schemaRegistryAdapter := repositories.NewSchemaRegistryAdapter(schemaRegistryClient)
	schemaService := services.NewSchemaService(schemaRegistryAdapter)
	handlers.NewHttpHandler(schemaService)
}
