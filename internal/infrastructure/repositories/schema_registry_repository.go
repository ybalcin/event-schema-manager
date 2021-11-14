package repositories

import (
	"../../../pkg/schemaregistry"
)

type (
	schemaRegistryRepository struct {
		schemaRegistryClient schemaregistry.IClient
	}
)

func NewSchemaRegistryRepository(schemaRegistryClient schemaregistry.IClient) *schemaRegistryRepository {
	return &schemaRegistryRepository{
		schemaRegistryClient: schemaRegistryClient,
	}
}

func (c *schemaRegistryRepository) Add(subject string, schema string) (int, error) {
	return c.schemaRegistryClient.RegisterNewSchema(subject, schema)
}
