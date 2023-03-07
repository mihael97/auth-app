package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mihael97/auth-proxy/src/util"
	"gitlab.com/mihael97/Go-utility/src/web"
	"gitlab.com/mihael97/Go-utility/src/web/routes"
)

var proxyControllerImpl *proxyController

type proxyController struct {
}

func (*proxyController) GetBasePath() string {
	return "/"
}

func (p *proxyController) GetRoutes() map[routes.Route]func(ctx *gin.Context) {
	return map[routes.Route]func(ctx *gin.Context){
		routes.CreateRoute("/*route", web.ALL, true): p.proxyRequests,
	}
}

func (p *proxyController) getRemoteUrl(ctx *gin.Context) (*url.URL, error) {
	path := ctx.Request.URL.Path
	if len(path) == 0 {
		return nil, fmt.Errorf("path is empty")
	}
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "api" {
		return nil, fmt.Errorf("path doesn't start with /api")
	}
	appName := parts[2]
	config := util.GetConfig()
	port, exists := config.ProxyServers[appName]
	if !exists {
		return nil, fmt.Errorf(fmt.Sprintf("app %s doesn't exist", appName))
	}
	newUrl := fmt.Sprintf("http://localhost:%d%s", port, path)
	log.Println("New url is ", newUrl)
	return url.Parse(newUrl)
}

func (p *proxyController) proxyRequests(ctx *gin.Context) {
	remote, err := p.getRemoteUrl(ctx)
	if err != nil {
		log.Printf("Error during creating remote: %v\n", err)
		web.WriteError(err, ctx)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = ctx.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = remote.Path
	}

	proxy.ServeHTTP(ctx.Writer, ctx.Request)
	if ctx.Writer.Status() == http.StatusBadGateway {
		web.WriteErrorMessage("error during proxying", ctx)
	}
}

func GetProxyController() routes.RoutesController {
	if proxyControllerImpl == nil {
		proxyControllerImpl = &proxyController{}
	}
	return proxyControllerImpl
}
