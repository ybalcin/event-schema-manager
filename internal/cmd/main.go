package main

import (
	"github.com/ybalcin/event-schema-manager/internal/cmd/http"
	"github.com/ybalcin/event-schema-manager/internal/shared/config"
)

func main() {
	cfg := config.LoadConfig()

	http.StartServer(cfg)
}
