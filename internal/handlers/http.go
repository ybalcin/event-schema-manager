package handlers

import (
	"../core/application/ports"
)

type (
	httpHandler struct {
		schemaService ports.ISchemaService
	}
)

func NewHttpHandler(schemaService ports.ISchemaService) *httpHandler {
	return &httpHandler{
		schemaService: schemaService,
	}
}
