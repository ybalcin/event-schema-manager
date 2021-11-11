package services

import (
	"../ports"
)

type (
	schemaService struct {
		schemaRepository ports.ISchemaRepository
	}
)

func NewSchemaService(schemaRepository ports.ISchemaRepository) *schemaService {
	return &schemaService{
		schemaRepository: schemaRepository,
	}
}
