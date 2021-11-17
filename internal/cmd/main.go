package main

import (
	"fmt"

	"github.com/ybalcin/event-schema-manager/internal/cmd/http"
	"github.com/ybalcin/event-schema-manager/internal/config"
)

func main() {
	cfg := config.LoadConfig()

	http.Start(cfg)

	fmt.Println("test4")
}
