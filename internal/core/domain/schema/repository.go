package schema

type (
	Repository interface {
		Add(subject string, schema string) (int, error)
	}
)
