package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mihael97/auth-proxy/src/util"
	goUtil "gitlab.com/mihael97/Go-utility/src/util"
	"gitlab.com/mihael97/Go-utility/src/web"
	middlewares "gitlab.com/mihael97/Go-utility/src/web/middlewares/security"
	"gitlab.com/mihael97/Go-utility/src/web/routes"
	"gitlab.com/mihael97/Go-utility/src/web/security/jwt"
)

const UsernameHeader = "X-MACUKA-USERNAME"
const RolesHeader = "X-MACUKA-ROLES"

func InitializeRoutes(engine *gin.Engine) {
	controllers := []routes.RoutesController{
		GetProxyController(),
	}

	engine.Use(middlewares.CORSMiddleware())
	engine.Use(JwtMiddleware())

	log.Print("Adding controller routes")
	routes.AddControllerRoutesWithFilter(false, engine, controllers)
}

func JwtMiddleware() func(ctx *gin.Context) {
	secret := util.GetConfig().Security.Secret
	return func(ctx *gin.Context) {
		appName, err := GetAppName(ctx)
		if err != nil {
			log.Println(err)
			ctx.Abort()
			return
		}
		appConfig, exists := util.GetConfig().ProxyServers[*appName]
		if !exists {
			web.ParseToJson(
				gin.H{"message": "app not found"},
				ctx,
				http.StatusNotFound,
			)
			ctx.Abort()
			return
		}

		route := ctx.FullPath()
		if _, exists := appConfig.SecuredRoutes[route]; exists {
			jwt.CheckSecurityToken(ctx, *secret)
		} else {
			log.Printf("Route %s is not secured\n", route)
			ctx.Next()
		}
	}
}

// CheckIfEligible Checks if user if eligible to access the endpoint with its roles
func CheckIfEligible(ctx *gin.Context) bool {
	appName, err := GetAppName(ctx)
	if err != nil {
		return false
	}
	config := util.GetConfig()
	appConfig, exists := config.ProxyServers[*appName]
	if !exists {
		ctx.Status(http.StatusNotFound)
		return false
	}
	route := ctx.Request.URL.Path
	appPath := ctx.FullPath()
	securedMethod, exists := appConfig.SecuredRoutes[appPath]
	if !exists {
		log.Printf("Route %s is not secured\n", route)
	} else {
		roles, exists := securedMethod[ctx.Request.Method]
		if !exists {
			log.Print("No additional roles check")
		} else {
			rolesHeader := ctx.Request.Header.Get(RolesHeader)
			if len(rolesHeader) == 0 {
				web.ParseToJson(
					gin.H{"message": "header not found"},
					ctx,
					http.StatusUnauthorized,
				)
				ctx.Abort()
				return false
			}
			if !goUtil.ContainsAny(roles, strings.Split(rolesHeader, ",")) {
				web.ParseToJson(
					gin.H{"message": "user doesn't have any eligible role for enter"},
					ctx,
					http.StatusForbidden,
				)
				ctx.Abort()
				return false
			}
		}
	}
	return true
}

func GetAppName(ctx *gin.Context) (*string, error) {
	path := ctx.Request.URL.Path
	if len(path) == 0 {
		return nil, fmt.Errorf("path is empty")
	}
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "api" {
		return nil, fmt.Errorf("path doesn't start with /api")
	}
	return goUtil.GetPointer(parts[2]), nil
}
