package security

import (
	"github.com/gin-gonic/gin"
	"github.com/mihael97/auth-proxy/src/util"
	goUtil "gitlab.com/mihael97/Go-utility/src/util"
	"gitlab.com/mihael97/Go-utility/src/web"
	"log"
	"net/http"
	"strings"
)

const UsernameHeader = "X-MACUKA-USERNAME"
const RolesHeader = "X-MACUKA-ROLES"
const IdHeader = "X-MACUKA-ID"
const ExpiresAtHeader = "X-MACUKA-EXPIRES-AT"

// CheckIfEligible Checks if user is eligible to access the endpoint with its roles
func CheckIfEligible(ctx *gin.Context) bool {
	appName, err := util.GetAppName(ctx)
	if err != nil {
		return false
	}
	config := util.GetConfig()
	appConfig, exists := config.ProxyServers[*appName]
	if !exists {
		ctx.Status(http.StatusNotFound)
		return false
	} else if len(appConfig.SecuredRoutes) == 0 {
		return true
	}

	appPath := ctx.Request.URL.Path
	securedMethod := appConfig.GetSecuredMethods(appPath)
	roles, exist := securedMethod[ctx.Request.Method]
	if !exist {
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
	return true
}
