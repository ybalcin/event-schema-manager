package adapters

import (
	"github.com/ybalcin/event-schema-manager/pkg/schemaregistry"
)

type (
	schemaRegistryRepository struct {
		schemaRegistryClient schemaregistry.Client
	}
)

func NewSchemaRegistryRepository(schemaRegistryClient schemaregistry.Client) *schemaRegistryRepository {
	return &schemaRegistryRepository{
		schemaRegistryClient: schemaRegistryClient,
	}
}

func (c *schemaRegistryRepository) Add(subject string, schema string) (int, error) {
	return c.schemaRegistryClient.RegisterNewSchema(subject, schema)
}
