package httpserver

import (
	"encoding/json"
	"os"
	"sync"
)

type (
	httpServerConfig struct {
		SchemaRegistryUrl string
	}
)

var config httpServerConfig

func loadConfig() {
	var once sync.Once
	once.Do(readConfig)
}

func readConfig() {
	var err error
	var file *os.File
	if file, err = os.Open("internal/cmd/httpserver/config.json"); err != nil {
		panic(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	config = httpServerConfig{}
	if err = decoder.Decode(&config); err != nil {
		panic(err)
	}
}
