package services

import (
	"github.com/ybalcin/event-schema-manager/internal/domain/schema"
)

type (
	ISchemaService interface {
		Add(subject string, schema string) error
	}
)

type (
	schemaService struct {
		schemaRepository schema.ISchemaRepository
	}
)

func NewSchemaService(schemaRepository schema.ISchemaRepository) ISchemaService {
	return &schemaService{
		schemaRepository: schemaRepository,
	}
}

func (s *schemaService) Add(subject string, schema string) error {
	_, err := s.schemaRepository.Add(subject, schema)
	// handle version, do something useful :)
	return err
}
