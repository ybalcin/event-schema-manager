package http

import (
	"../../../core/application/ports"
	"../../../shared/config"
	"fmt"
)

func Start(cfg *config.AppConfig) {
	if cfg == nil {
		panic("AppConfig is nil!")
	}

	ports.NewHttpServer(cfg)

	fmt.Println("HttpServer running!")
}
