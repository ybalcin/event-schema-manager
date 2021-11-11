package ports

type (
	ISchemaService interface {
		Add(subject string, schema string) error
	}
)

