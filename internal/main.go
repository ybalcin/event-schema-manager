package main

import (
	"./cmd/server/http"
	"./shared/config"
)

func main() {
	cfg := config.LoadConfig()

	http.Start(cfg)
}
