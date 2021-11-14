package application

import (
	"../../../pkg/schemaregistry"
	"../../infrastructure/repositories"
	"../../shared/config"
	. "./services"
)

type (
	Application struct {
		schemaService ISchemaService
	}
)

func NewApplication(cfg *config.AppConfig) *Application {
	schemaRegistryClient, err := schemaregistry.NewClient(cfg.SchemaRegistryUrl)
	if err != nil {
		panic(err)
	}

	schemaRegistryAdapter := repositories.NewSchemaRegistryRepository(schemaRegistryClient)
	schemaService := NewSchemaService(schemaRegistryAdapter)
	return &Application{schemaService: schemaService}
}
