package ports

type (
	ISchemaRepository interface {
		Add(subject string, schema string) (int, error)
	}
)
