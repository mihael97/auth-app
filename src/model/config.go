package model

type Config struct {
	Port         *string        `yaml:"port"`
	ProxyServers map[string]int `yaml:"proxyServers"`
}
