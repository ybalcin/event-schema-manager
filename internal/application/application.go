package application

import (
	"github.com/ybalcin/event-schema-manager/internal/adapters"
	"github.com/ybalcin/event-schema-manager/internal/config"
	"github.com/ybalcin/event-schema-manager/internal/services"
	"github.com/ybalcin/event-schema-manager/pkg/schemaregistry"
)

type (
	Application struct {
		schemaService services.ISchemaService
	}
)

func NewApplication(cfg *config.AppConfig) *Application {
	schemaRegistryClient, err := schemaregistry.NewClient(cfg.SchemaRegistryUrl)
	if err != nil {
		panic(err)
	}

	schemaRegistryAdapter := adapters.NewSchemaRegistryRepository(schemaRegistryClient)
	schemaService := services.NewSchemaService(schemaRegistryAdapter)
	return &Application{schemaService: schemaService}
}
