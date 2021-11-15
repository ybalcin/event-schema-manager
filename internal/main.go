package main

import (
	"fmt"

	"./cmd/server/http"
	"./shared/config"
)

func main() {
	cfg := config.LoadConfig()

	http.Start(cfg)

	fmt.Println("test4")
}
