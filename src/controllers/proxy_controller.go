package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mihael97/auth-proxy/src/services"
	"github.com/mihael97/auth-proxy/src/util"
	exceptionUtil "gitlab.com/mihael97/Go-utility/src/util"
	"gitlab.com/mihael97/Go-utility/src/web"
	"gitlab.com/mihael97/Go-utility/src/web/routes"
	"gitlab.com/mihael97/Go-utility/src/web/security/jwt"
)

const UsernameHeader = "X-MACUKA-USERNAME"
const RolesHeader = "X-MACUKA-ROLES"

var proxyControllerImpl *proxyController

type proxyController struct {
	routingTable map[string]*gin.Engine
	userService  services.UserService
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
	if len(parts) < 2 || parts[1] != "api" {
		return nil, false, fmt.Errorf("path doesn't start with /api")
	}
	appName := parts[2]
	config := util.GetConfig()
	port, exists := config.ProxyServers[appName]
	if !exists {
		log.Printf("URL %s not found\n", appName)
		return nil, false, nil
	}
	path = strings.ReplaceAll(path, "/"+appName, "")
	newUrl := fmt.Sprintf("http://localhost:%d%s/", port, path)
	log.Println("New url is ", newUrl)
	newPath, err := url.Parse(newUrl)
	return newPath, true, err
}

func (p *proxyController) proxyRequests(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	searchPath := strings.Join(strings.Split(path, "/")[0:3], "/")
	if router, exits := p.routingTable[searchPath]; exits {
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
		web.WriteErrorMessage("error during modifiying headers", ctx)
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
	p.modifyHeaders(ctx, true)
	if ctx.Writer.Status() == http.StatusBadGateway {
		web.WriteErrorMessage("error during proxying", ctx)
	}
}

// Add username and roles header
func (p *proxyController) modifyHeaders(ctx *gin.Context, remove ...bool) bool {
	if len(remove) != 0 && remove[0] {
		ctx.Writer.Header().Del(UsernameHeader)
		ctx.Writer.Header().Del(RolesHeader)
	} else {
		username, err := jwt.GetUserNameFromToken(ctx, *util.GetConfig().Security.Secret)
		userData, _ := p.userService.GetUser(username)
		if err != nil {
			log.Println("Error during parsing token", err)
			return false
		}
		ctx.Request.Header.Add(UsernameHeader, username)
		ctx.Request.Header.Add(RolesHeader, strings.Join(userData.Roles, ","))
	}
	return true
}

func GetProxyController() routes.RoutesController {
	if proxyControllerImpl == nil {
		proxyControllerImpl = &proxyController{
			map[string]*gin.Engine{"/api/login": GetLoginController(), "/api/users": GetUserController()},
			services.GetUserService(),
		}
	}
	return proxyControllerImpl
}
