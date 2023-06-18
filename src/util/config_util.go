package util

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/mihael97/auth-proxy/src/model"
	"gitlab.com/mihael97/Go-utility/src/env"
	"gitlab.com/mihael97/Go-utility/src/web/routes"
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

	config.ProxyServers = make(map[string]struct {
		Port            *int
		SecuredRoutes   map[string]map[string][]string
		UnsecuredRoutes map[string][]string "yaml:\"unsecuredRoutes\""
	}, 0)
	for _, backend := range config.Backends {
		client := http.DefaultClient
		url := fmt.Sprintf("%s/api/routes", backend)
		response, err := client.Get(url)
		if err != nil {
			log.Panic(err)
			return
		}
		var allRoutes map[string][]routes.Route
		if err := json.NewDecoder(response.Body).Decode(&allRoutes); err != nil {
			log.Panic(err)
		}

		portParts := strings.Split(url, "/")
		port, _ := strconv.ParseInt(portParts[len(portParts)-1], 10, 64)
		portInt := int(port)
		for routePath, routes := range allRoutes {
			unsecuredRoutes := make(map[string][]string, 0)

			for _, route := range routes {
				if !route.Secured {
					value, exist := unsecuredRoutes[route.URL]
					if !exist {
						value = make([]string, 0)
					}
					value = append(value, string(route.Type))
					unsecuredRoutes[route.URL] = value
				}
			}

			config.ProxyServers[routePath] = struct {
				Port            *int
				SecuredRoutes   map[string]map[string][]string
				UnsecuredRoutes map[string][]string "yaml:\"unsecuredRoutes\""
			}{
				Port:            &portInt,
				UnsecuredRoutes: unsecuredRoutes,
			}
		}

	}
}
