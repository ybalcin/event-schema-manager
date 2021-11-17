package http

import (
	"fmt"
	"github.com/ybalcin/event-schema-manager/internal/core/application/ports"
	"github.com/ybalcin/event-schema-manager/internal/shared/config"
)

func StartServer(cfg *config.AppConfig) {
	if cfg == nil {
		panic("AppConfig is nil!")
	}

	ports.NewHttpServer(cfg)

	fmt.Println("HttpServer running!")
}
