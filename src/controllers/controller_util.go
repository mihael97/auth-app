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
	"gitlab.com/mihael97/Go-utility/src/web/routes"
)

const UsernameHeader = "X-MACUKA-USERNAME"
const RolesHeader = "X-MACUKA-ROLES"

func InitializeRoutes(engine *gin.Engine) {
	controllers := []routes.RoutesController{
		GetProxyController(),
	}

	log.Print("Adding controller routes")
	routes.AddControllerRoutesWithFilter(false, engine, controllers)
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
	if goUtil.Contains(appPath, appConfig.UnsecuredRoutes...) {
		log.Printf("Route %s is not secured\n", route)
	} else {
		RolesHeader := ctx.Request.Header.Get(UsernameHeader)
		if len(RolesHeader) == 0 {
			web.ParseToJson(
				gin.H{"message": "header not found"},
				ctx,
				http.StatusUnauthorized,
			)
			ctx.Abort()
			return false
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
