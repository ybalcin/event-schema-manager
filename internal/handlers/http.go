package handlers

import (
	"../core/ports"
)

type (
	HttpHandler struct {
		schemaService ports.ISchemaService
	}
)

func NewHttpHandler(schemaService ports.ISchemaService) *HttpHandler {
	return &HttpHandler{
		schemaService: schemaService,
	}
}
