package application

import (
	"github.com/ybalcin/event-schema-manager/internal/core/application/services"
	"github.com/ybalcin/event-schema-manager/internal/infrastructure/adapters"
	"github.com/ybalcin/event-schema-manager/internal/shared/config"
	"github.com/ybalcin/event-schema-manager/pkg/schemaregistry"
)

type (
	Application struct {
		schemaService services.SchemaService
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
