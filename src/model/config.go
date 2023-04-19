package model

type Config struct {
	Security struct {
		Secret         *string `yaml:"secret"`
		ValidityPeriod *uint64 `yaml:"validityPerion"`
		HeaderName     *string `yaml:"headerName"`
	}
	Port         *string `yaml:"port"`
	ProxyServers map[string]struct {
		Port            *int
		SecuredRoutes   map[string]map[string][]string
		UnsecuredRoutes map[string][]string `yaml:"unsecuredRoutes"`
	} `yaml:"proxyServers"`
}
