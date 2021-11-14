package main

import (
	"../shared/config"
	"./server/http"
)

func main() {
	cfg := config.LoadConfig()

	http.Start(cfg)
}
