package application

import (
	. "../services"
	"github.com/ybalcin/event-schema-manager/internal/repositories"
	"github.com/ybalcin/event-schema-manager/internal/shared/config"
	"github.com/ybalcin/event-schema-manager/pkg/schemaregistry"
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
