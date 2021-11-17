package config

import (
	"encoding/json"
	"os"
	"sync"
)

type (
	AppConfig struct {
		SchemaRegistryUrl string
	}
)

var config AppConfig

func LoadConfig() *AppConfig {
	var once sync.Once
	once.Do(readFromJson)

	return &config
}

func Config() *AppConfig {
	return &config
}

func readFromJson() {
	var err error
	var file *os.File
	if file, err = os.Open("../config/config.json"); err != nil {
		panic(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	config = AppConfig{}
	if err = decoder.Decode(&config); err != nil {
		panic(err)
	}
}
