package services

import (
	. "../domain/schema"
)

type (
	ISchemaService interface {
		Add(subject string, schema string) error
	}
)

type (
	schemaService struct {
		schemaRepository ISchemaRepository
	}
)

func NewSchemaService(schemaRepository ISchemaRepository) ISchemaService {
	return &schemaService{
		schemaRepository: schemaRepository,
	}
}

func (s *schemaService) Add(subject string, schema string) error {
	_, err := s.schemaRepository.Add(subject, schema)
	// handle version, do something useful :)
	return err
}
