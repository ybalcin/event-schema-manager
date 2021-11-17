package ports

import (
	"../config"
	"./.."
)

// HttpServer input http port
type (
	HttpServer struct {
		app *application.Application
	}
)

func NewHttpServer(cfg *config.AppConfig) *HttpServer {
	app := application.NewApplication(cfg)
	return &HttpServer{app: app}
}
