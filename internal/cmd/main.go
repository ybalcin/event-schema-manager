package main

import (
	"fmt"

	"github.com/ybalcin/event-schema-manager/internal/config"
	"github.com/ybalcin/event-schema-manager/internal/http"
)

func main() {
	cfg := config.LoadConfig()

	http.Start(cfg)

	fmt.Println("test4")
}
