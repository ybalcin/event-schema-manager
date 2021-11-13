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

func (s *schemaService) Add(subject string, schema string) error {
	_, err := s.schemaRepository.Add(subject, schema)
	// handle version, do something useful :)
	return err
}
