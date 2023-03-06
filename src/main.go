package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mihael97/auth-proxy/src/controllers"
	"github.com/mihael97/auth-proxy/src/model"
	"gitlab.com/mihael97/Go-utility/src/env"
	"gopkg.in/yaml.v3"
)

func readConfigFile() {
	configPath := env.GetEnvVariable("PROXY_CONFIG", "./config.yaml")
	if len(configPath) == 0 {
		log.Panic("Config file should exist\n")
	}
	file, err := os.ReadFile(configPath)
	if err != nil {
		log.Panic("Error during reading file\n", err)
	}
	var config model.Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Panic("Yaml error", err)
	}
}

func main() {
	port := env.GetEnvVariable("HTTP_PORT", "8080")
	log.Printf("Auth proxy starting at port %s", port)

	readConfigFile()
	router := gin.Default()
	controllers.InitializeRoutes(router)

	log.Panic(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
