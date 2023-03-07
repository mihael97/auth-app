package util

import (
	"log"
	"os"

	"github.com/mihael97/auth-proxy/src/model"
	"gitlab.com/mihael97/Go-utility/src/env"
	"gopkg.in/yaml.v2"
)

var config *model.Config

func GetConfig() *model.Config {
	if config == nil {
		readConfigFile()
	}
	return config
}

func readConfigFile() {
	configPath := env.GetEnvVariable("PROXY_CONFIG", "./config.yaml")
	if len(configPath) == 0 {
		log.Panic("Config file should exist\n")
	}
	file, err := os.ReadFile(configPath)
	if err != nil {
		log.Panic("Error during reading file\n", err)
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Panic("Yaml error", err)
	}
}
