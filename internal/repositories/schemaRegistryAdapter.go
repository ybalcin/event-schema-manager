package repositories

type (
	schemaRegistryAdapter struct {
	}
)

func NewSchemaRegistryAdapter() *schemaRegistryAdapter {
	return &schemaRegistryAdapter{}
}
