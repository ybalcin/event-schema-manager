package ports

import (
	"github.com/ybalcin/event-schema-manager/internal/app"
	"github.com/ybalcin/event-schema-manager/internal/shared/config"
)

// HttpServer input http port
type (
	HttpServer struct {
		app *app.Application
	}
)

func NewHttpServer(cfg *config.AppConfig) *HttpServer {
	app := app.NewApplication(cfg)
	return &HttpServer{app: app}
}
