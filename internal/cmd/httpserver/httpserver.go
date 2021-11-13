package httpserver

import (
	"../../../internal/handlers"
	"../../../pkg/schemaregistry"
	"../../core/application/services"
	"../../infrastructure/repositories"
	"fmt"
)

func Start() {
	loadConfig()

	schemaRegistryClient, err := schemaregistry.NewClient(config.SchemaRegistryUrl)
	if err != nil {
		panic(err)
	}

	schemaRegistryAdapter := repositories.NewSchemaRegistryAdapter(schemaRegistryClient)
	schemaService := services.NewSchemaService(schemaRegistryAdapter)
	handlers.NewHttpHandler(schemaService)

	fmt.Println("HttpServer running!")
}
