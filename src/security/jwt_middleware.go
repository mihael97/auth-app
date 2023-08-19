package security

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mihael97/auth-proxy/src/util"
	goUtil "gitlab.com/mihael97/Go-utility/src/util"
	"gitlab.com/mihael97/Go-utility/src/web/security/jwt"
	"log"
	"regexp"
	"strings"
)

const SwaggerResources = "/swagger/*"

func JwtMiddleware() func(ctx *gin.Context) {
	secret := util.GetConfig().Security.Secret
	return func(ctx *gin.Context) {
		// remove public header
		ctx.Request.Header.Del("public")

		appName, err := util.GetAppName(ctx)
		if err != nil {
			log.Println(err)
			ctx.Abort()
			return
		}
		appConfig, exists := util.GetConfig().ProxyServers[*appName]

		if exists {
			route := strings.ReplaceAll(ctx.Request.URL.Path, fmt.Sprintf("/%s", *appName), "")
			if route == "/api/routes" {
				ctx.Request.Header.Add("public", "true")
				log.Println("Accessing routes")
				ctx.Next()
				return
			}

			route = strings.TrimPrefix(route, "/api")
			unsecuredRouteMethods, exists := appConfig.UnsecuredRoutes[route]

			if ok, _ := regexp.Match(SwaggerResources, []byte(route)); ok || exists {
				if len(unsecuredRouteMethods) == 0 || goUtil.Contains(ctx.Request.Method, unsecuredRouteMethods...) {
					log.Printf("Route %s is not secured\n", route)
					ctx.Request.Header.Add("public", "true")
					ctx.Next()
					return
				}
			}
			if exists := appConfig.IsSecured(route); len(appConfig.SecuredRoutes) == 0 || exists {
				jwt.CheckSecurityToken(ctx, *secret)
			} else {
				log.Printf("Route %s is not secured\n", route)
				ctx.Next()
			}
		}
	}
}
