package main

import (
	"encoding/json"
	"os"
	"sync"
)

type (
	appConfig struct {
		schemaRegistryUrl string
	}
)

var config appConfig

func init() {
	var once sync.Once
	once.Do(readConfig)
}

func readConfig() {
	var err error
	var file *os.File
	if file, err = os.Open("/config.json"); err != nil {
		panic(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	config = appConfig{}
	if err = decoder.Decode(&config); err != nil {
		panic(err)
	}
}
