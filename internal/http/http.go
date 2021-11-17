package http

import (
	"fmt"
	"github.com/ybalcin/event-schema-manager/internal/config"
	"github.com/ybalcin/event-schema-manager/internal/ports"
)

func Start(cfg *config.AppConfig) {
	if cfg == nil {
		panic("AppConfig is nil!")
	}

	ports.NewHttpServer(cfg)

	fmt.Println("HttpServer running!")
}
