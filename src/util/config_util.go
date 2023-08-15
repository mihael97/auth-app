package util

import (
	"encoding/json"
	"fmt"
	"gitlab.com/mihael97/Go-utility/src/util"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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

	appendBackends(config)
	appendEmailConfig(config)
	appendPasswordRecovery(config)
}

func appendPasswordRecovery(c *model.Config) {
	if c.PasswordRecovery.Url == nil {
		c.PasswordRecovery.Url = util.GetPointer("http://localhost:3000")
	}
}

func appendEmailConfig(c *model.Config) {
	if c.EmailConfig.From == nil {
		c.EmailConfig.From = util.GetPointer("mihaelmacuka2@gmail.com")
	}
	if c.EmailConfig.ServerConfig.Secret == nil {
		log.Panic("Email secret should be provided")
	}
	if c.EmailConfig.ServerConfig.Host == nil {
		c.EmailConfig.ServerConfig.Host = util.GetPointer("smtp.gmail.com")
	}
	if c.EmailConfig.ServerConfig.Port == nil {
		c.EmailConfig.ServerConfig.Port = util.GetPointer("587")
	}
}

func appendBackends(config *model.Config) {
	config.ProxyServers = make(map[string]model.ProxyServer, 0)
	for name, backendConfig := range config.Backends {
		if !backendConfig.IsEnabled() {
			log.Printf("%s is disabled", backendConfig.Url)
			continue
		}

		go parseBackendConfig(config, backendConfig, name)
	}
}

func parseBackendConfig(config *model.Config, backendConfig model.BackendServerConfig, name string) {
	url, allRoutes := fetchBackendRoutes(backendConfig)

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

func fetchBackendRoutes(backendConfig model.BackendServerConfig) (string, map[string][]routes.Route) {
	client := http.DefaultClient
	url := fmt.Sprintf("%s/routes", backendConfig.Url)
	var response *http.Response

	retrySeconds, _ := strconv.ParseInt(env.GetEnvVariable("BACKEND_RETRY_PERIOD", "5"), 10, 64)
	retryMax, _ := strconv.ParseInt(env.GetEnvVariable("BACKEND_RETRY_MAX", "5"), 10, 64)

	retryCount := 0
	for {
		if retryCount == int(retryMax) {
			log.Panicf("After %d retries, cannot connect to %s", retryCount, url)
		}
		retryCount += 1

		var err error
		response, err = client.Get(url)

		if response != nil && response.StatusCode == http.StatusOK {
			break
		} else if err != nil {
			log.Printf("Error. Retry: %d, URL: %s\nError: %s\n", retryCount, url, err.Error())
		}
		responseStatus := "N/A"
		if response != nil {
			responseStatus = response.Status
		}
		log.Printf("Retry: %d, status: %s. Retry after %d seconds", retryCount, responseStatus, retrySeconds)
		sleep(int(retrySeconds))
	}

	var allRoutes map[string][]routes.Route
	if err := json.NewDecoder(response.Body).Decode(&allRoutes); err != nil {
		log.Panic(err)
	}

	return url, allRoutes
}

func sleep(seconds int) {
	for i := 0; i < seconds; i++ {
		time.Sleep(1 * time.Second)
	}
}
