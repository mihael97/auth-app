package model

type Config struct {
	Security struct {
		Secret         *string `yaml:"secret"`
		ValidityPeriod *uint64 `yaml:"validityPerion"`
		HeaderName     *string `yaml:"headerName"`
	}
	Port         *string        `yaml:"port"`
	ProxyServers map[string]int `yaml:"proxyServers"`
}
