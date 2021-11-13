package httpserver

import (
	"../../../internal/handlers"
	"../../../internal/repositories"
	"../../../pkg/schemaregistry"
	"../../core/services"
)

func Run() {
	loadConfig()

	schemaRegistryClient, err := schemaregistry.NewClient(config.SchemaRegistryUrl)
	if err != nil {
		panic(err)
	}

	schemaRegistryAdapter := repositories.NewSchemaRegistryAdapter(schemaRegistryClient)
	schemaService := services.NewSchemaService(schemaRegistryAdapter)
	handlers.NewHttpHandler(schemaService)
}
