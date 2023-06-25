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
		config = &model.Config{}
		readConfigFile(config)
	}
	return config
}

func readConfigFile(config *model.Config) {
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

	config.ProxyServers = make(map[string]model.ProxyServer, 0)
	for name, path := range config.Backends {
		client := http.DefaultClient
		url := fmt.Sprintf("%s/routes", path)
		response, err := client.Get(url)
		if err != nil {
			log.Panic(err)
		}

		if response.StatusCode != http.StatusOK {
			log.Panic("Status is ", response.StatusCode)
		}

		var allRoutes map[string][]routes.Route
		if err := json.NewDecoder(response.Body).Decode(&allRoutes); err != nil {
			log.Panic(err)
		}

		portParts := strings.Split(url, "/")
		portParts = strings.Split(portParts[len(portParts)-2], ":")
		port, _ := strconv.ParseInt(portParts[1], 10, 64)
		portInt := int(port)
		unsecuredRoutes := make(map[string][]string, 0)
		securedRoutes := make(map[string]map[string][]string, 0)
		for routePath, routes := range allRoutes {

			for _, route := range routes {
				urlPath := fmt.Sprintf("%s%s", routePath, route.URL)
				if route.Secured {
					value, exist := securedRoutes[urlPath]
					if !exist {
						value = make(map[string][]string, 0)
					}
					subValue, exist := value[route.Type.String()]
					if !exist {
						subValue = make([]string, 0)
					}
					if route.Roles != nil {
						subValue = append(subValue, *route.Roles...)
					}
					value[route.Type.String()] = subValue
					securedRoutes[urlPath] = value
				} else {
					value, exist := unsecuredRoutes[urlPath]
					if !exist {
						value = make([]string, 0)
					}
					value = append(value, string(route.Type))
					unsecuredRoutes[urlPath] = value
				}
			}
		}

		config.ProxyServers[name] = model.ProxyServer{
			Port:            &portInt,
			UnsecuredRoutes: unsecuredRoutes,
			SecuredRoutes:   securedRoutes,
		}
	}
}
