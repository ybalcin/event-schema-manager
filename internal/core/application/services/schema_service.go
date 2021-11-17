package services

import (
	"github.com/ybalcin/event-schema-manager/internal/core/domain/schema"
)

type (
	SchemaService interface {
		Add(subject string, schema string) error
	}
)

type (
	schemaService struct {
		repository schema.Repository
	}
)

func NewSchemaService(repository schema.Repository) SchemaService {
	return &schemaService{
		repository: repository,
	}
}

func (s *schemaService) Add(subject string, schema string) error {
	_, err := s.repository.Add(subject, schema)
	// handle version, do something useful :)
	return err
}
