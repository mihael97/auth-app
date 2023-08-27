package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mihael97/auth-proxy/docs"
	"github.com/mihael97/auth-proxy/src/security"
	config "github.com/mihael97/auth-proxy/src/util"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gitlab.com/mihael97/Go-utility/src/web"
	middlewares "gitlab.com/mihael97/Go-utility/src/web/middlewares/security"
	"gitlab.com/mihael97/Go-utility/src/web/routes"
	"gitlab.com/mihael97/Go-utility/src/web/security/jwt"
	"log"
	"strconv"
	"strings"
)

func InitializeRoutes(engine *gin.Engine) {
	controllers := []routes.RoutesController{
		GetProxyController(),
	}

	engine.Use(middlewares.CORSMiddleware())
	engine.Use(security.JwtMiddleware())

	log.Print("Adding controller routes")
	routes.AddControllerRoutesWithFilter(false, engine, controllers)
}

func InitSwagger() *gin.Engine {
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Title = fmt.Sprintf("Auth backend")
	swaggerPath := "/api/swagger/*any"
	engine := gin.New()
	engine.GET(swaggerPath, ginSwagger.WrapHandler(swaggerFiles.Handler))
	return engine
}

func appendExpiresAt(ctx *gin.Context) {
	var remove = false
	if len(ctx.Request.Header.Get(security.AuthorizationHeader)) == 0 {
		remove = true
		authorizationHeader := ctx.Writer.Header().Get(security.AuthorizationHeader)
		if !strings.HasPrefix("Bearer", authorizationHeader) {
			authorizationHeader = fmt.Sprintf("Bearer %s", authorizationHeader)
		}
		ctx.Request.Header.Add(security.AuthorizationHeader, authorizationHeader)
	}
	maker, err := jwt.NewJwtMaker(*config.GetConfig().Security.Secret)
	if err != nil {
		web.WriteError(err, ctx)
		return
	}
	payload, err := maker.VerifyToken(ctx)
	if err != nil {
		web.WriteError(err, ctx)
		return
	}
	expiresAtTime := payload.ExpiredAt.Unix()
	ctx.Writer.Header().Add(security.ExpiresAtHeader, strconv.FormatInt(expiresAtTime, 10))
	if remove {
		ctx.Request.Header.Del(security.AuthorizationHeader)
	}
}
