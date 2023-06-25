package model

type SecurityConfig struct {
	Secret         *string `yaml:"secret"`
	ValidityPeriod *uint64 `yaml:"validityPerion"`
	HeaderName     *string `yaml:"headerName"`
}

type ProxyServer struct {
	Port            *int                           `yaml:"port"`
	SecuredRoutes   map[string]map[string][]string `yaml:"securedRoutess"`
	UnsecuredRoutes map[string][]string            `yaml:"unsecuredRoutes"`
}

func (p ProxyServer) GetSecuredMethods(url string) map[string][]string {
	for key, value := range p.SecuredRoutes {
		if p.isMatch(key, url) {
			return value
		}
	}

	return nil
}

func (p ProxyServer) isMatch(key, url string) bool {
	keyIndex := 0
	urlIndex := 0
	for keyIndex < len(key) && urlIndex < len(url) {
		character := key[keyIndex]
		skip := false
		if character == ':' && key[keyIndex-1] == '/' {
			skip = true
			for keyIndex < len(key) && key[keyIndex] != '/' {
				keyIndex += 1
			}
			for keyIndex < len(url) && key[urlIndex] != '/' {
				urlIndex += 1
			}
		}
		if !skip && key[keyIndex] != url[urlIndex] {
			return false
		}
		urlIndex += 1
		keyIndex += 1
	}

	if keyIndex == len(key)-1 && (urlIndex == len(url)-1 || url[urlIndex+1] == '?') {
		return true
	}

	return false
}

func (p ProxyServer) IsSecured(url string) bool {
	isSecured := false

	for route := range p.SecuredRoutes {
		if route == url || route == url+"/" {
			isSecured = true
			break
		}
	}

	return isSecured
}

type Config struct {
	Security     SecurityConfig         `yaml:"security"`
	Port         *string                `yaml:"port"`
	Backends     map[string]string      `yaml:"backends"`
	ProxyServers map[string]ProxyServer `yaml:"proxyServers"`
}
