package model

type Config struct {
	Security struct {
		Secret         *string `yaml:"secret"`
		ValidityPeriod *uint64 `yaml:"validityPerion"`
	}
	Port         *string        `yaml:"port"`
	ProxyServers map[string]int `yaml:"proxyServers"`
}
