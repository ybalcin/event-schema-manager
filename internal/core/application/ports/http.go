package ports

import (
	"github.com/ybalcin/event-schema-manager/internal/core/application"
	"github.com/ybalcin/event-schema-manager/internal/shared/config"
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
