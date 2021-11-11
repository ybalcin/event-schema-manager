package repositories

import (
	"../../pkg/schemaregistry"
)

type (
	schemaRegistryAdapter struct {
		schemaRegistryClient schemaregistry.IClient
	}
)

func NewSchemaRegistryAdapter(schemaRegistryClient schemaregistry.IClient) *schemaRegistryAdapter {
	return &schemaRegistryAdapter{
		schemaRegistryClient: schemaRegistryClient,
	}
}

func (c *schemaRegistryAdapter) Add(subject string, schema string) (int, error) {
	return c.schemaRegistryClient.RegisterNewSchema(subject, schema)
}