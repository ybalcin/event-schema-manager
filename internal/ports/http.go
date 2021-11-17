package ports

import (
	"github.com/ybalcin/event-schema-manager/internal/application"
	"github.com/ybalcin/event-schema-manager/internal/config"
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
