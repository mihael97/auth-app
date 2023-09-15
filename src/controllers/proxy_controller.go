package controllers

import (
	"fmt"
	"github.com/mihael97/auth-proxy/src/security"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mihael97/auth-proxy/src/services"
	"github.com/mihael97/auth-proxy/src/util"
	exceptionUtil "gitlab.com/mihael97/Go-utility/src/util"
	goUtil "gitlab.com/mihael97/Go-utility/src/util"
	"gitlab.com/mihael97/Go-utility/src/web"
	"gitlab.com/mihael97/Go-utility/src/web/routes"
	"gitlab.com/mihael97/Go-utility/src/web/security/jwt"
)

const ApiVersionQueryParam = "apiVersion"
const ApiVersionRegex = "[vV][0-9]+"
const ApiVersionHeader = "MACUKA_API_VERSION"
const ApiVersionV1 = "V1"

var proxyControllerImpl *proxyController

type proxyController struct {
	routingTable            map[string]*gin.Engine
	userService             services.UserService
	permittedUsersEndpoints []string
}

func (*proxyController) GetBasePath() string {
	return "/"
}

func (p *proxyController) GetRoutes() map[routes.Route]func(ctx *gin.Context) {
	return map[routes.Route]func(ctx *gin.Context){
		routes.CreateRoute("/*route", web.ALL, true): p.proxyRequests,
	}
}

func (p *proxyController) getRemoteUrl(ctx *gin.Context) (*url.URL, bool, error) {
	path := ctx.Request.URL.Path
	if len(path) == 0 {
		return nil, false, fmt.Errorf("path is empty")
	}
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "api" && parts[1] != "routes" {
		return nil, false, fmt.Errorf("path doesn't start with /api or /routes")
	}

	p.parseVersion(ctx, &path)

	appName, err := util.GetAppName(ctx)
	if err != nil {
		return nil, false, err
	}

	config := util.GetConfig()
	serverData, exists := config.ProxyServers[*appName]
	if !exists {
		log.Printf("URL %s not found\n", *appName)
		return nil, false, nil
	}
	path = strings.ReplaceAll(path, fmt.Sprintf("/%s", *appName), "")
	endPath := ""
	if len(parts) <= 4 && path != "/api/routes" {
		endPath = "/"
	}
	if path == "/api/routes" {
		path = strings.TrimPrefix(path, "/api")
	}
	newUrl := fmt.Sprintf("http://localhost:%d%s%s", *serverData.Port, path, endPath)
	log.Println("New url is ", newUrl)
	newPath, err := url.Parse(newUrl)
	return newPath, true, err
}

func (p *proxyController) proxyRequests(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	searchPath := strings.Join(strings.Split(path, "/")[0:3], "/")
	if router, exists := p.routingTable[searchPath]; exists {
		requestUri := ctx.Request.RequestURI
		if strings.HasPrefix(requestUri, "/api/users") && !p.isEndpointPermittedForAll(requestUri) {
			modifyHeadersStatus := p.modifyHeaders(ctx)
			if !modifyHeadersStatus {
				web.WriteErrorMessage("error during modifying headers", ctx)
				return
			}
			roles := strings.Split(ctx.Request.Header.Get(security.RolesHeader), ",")
			if !goUtil.Contains("ADMIN", roles...) {
				ctx.AbortWithStatus(http.StatusForbidden)
				return
			}
		}
		router.HandleContext(ctx)
		return
	}

	remote, found, err := p.getRemoteUrl(ctx)
	if err != nil {
		log.Printf("Error during creating remote: %v\n", err)
		web.WriteError(err, ctx)
		return
	} else if !found {
		errorResponse := exceptionUtil.NewException("path not found")
		web.ParseToJson(errorResponse, ctx, http.StatusNotFound)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	modifyHeadersStatus := p.modifyHeaders(ctx)
	if !modifyHeadersStatus {
		web.WriteErrorMessage("error during modifying headers", ctx)
		return
	}

	// check if eligible
	if !security.CheckIfEligible(ctx) {
		ctx.Abort()
		return
	}

	proxy.Director = func(req *http.Request) {
		req.Header = ctx.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = remote.Path
	}

	proxy.ServeHTTP(ctx.Writer, ctx.Request)
	p.modifyHeaders(ctx, false)
	if ctx.Writer.Status() == http.StatusBadGateway {
		web.WriteErrorMessage("error during proxying", ctx)
	}
}

// Add username and roles header
func (p *proxyController) modifyHeaders(ctx *gin.Context, add ...bool) bool {
	ctx.Writer.Header().Del(security.UsernameHeader)
	ctx.Writer.Header().Del(security.RolesHeader)
	ctx.Writer.Header().Del(security.IdHeader)
	ctx.Writer.Header().Del(security.ExpiresAtHeader)

	if len(add) == 0 || add[0] {
		if len(ctx.Request.Header.Get("public")) != 0 {
			return true
		}
		username, err := jwt.GetUserNameFromToken(ctx, *util.GetConfig().Security.Secret)
		if err != nil {
			log.Println("Error during parsing token", err)
			return false
		}
		userData, err := p.userService.GetUser(username)
		if err != nil {
			log.Println("Error during fetching user info", err)
			return false
		}
		ctx.Request.Header.Add(security.UsernameHeader, username)
		ctx.Request.Header.Add(security.IdHeader, userData.Id)
		ctx.Request.Header.Add(security.RolesHeader, strings.Join(userData.Roles, ","))
		appendExpiresAt(ctx)
	}
	return true
}

func (p *proxyController) isEndpointPermittedForAll(requestPath string) bool {
	for _, route := range p.permittedUsersEndpoints {
		if strings.HasPrefix(requestPath, route) {
			return true
		}
	}
	return false
}

func (p *proxyController) parseVersion(ctx *gin.Context, path *string) {
	apiVersion := ctx.Query(ApiVersionQueryParam)
	match, _ := regexp.MatchString(ApiVersionRegex, apiVersion)
	if !match {
		apiVersion = ApiVersionV1
	}
	ctx.Request.Header.Set(ApiVersionHeader, apiVersion)
}

func GetProxyController() routes.RoutesController {
	if proxyControllerImpl == nil {
		proxyControllerImpl = &proxyController{
			map[string]*gin.Engine{"/api/login": GetLoginController(), "/api/users": GetUserController(), "/api/swagger": InitSwagger()},
			services.GetUserService(),
			[]string{"/api/users/recovery", "/api/users/me"},
		}
	}
	return proxyControllerImpl
}
