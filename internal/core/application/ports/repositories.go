package ports

// ISchemaRepository output schema repository port
type (
	ISchemaRepository interface {
		Add(subject string, schema string) (int, error)
	}
)
