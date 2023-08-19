package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mihael97/auth-proxy/src/security"
	config "github.com/mihael97/auth-proxy/src/util"
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

func appendExpiresAt(ctx *gin.Context) {
	if len(ctx.Request.Header.Get("Authorization")) == 0 {
		authorizationHeader := ctx.Writer.Header().Get("Authorization")
		if !strings.HasPrefix("Bearer", authorizationHeader) {
			authorizationHeader = fmt.Sprintf("Bearer %s", authorizationHeader)
		}
		ctx.Request.Header.Add("Authorization", authorizationHeader)
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
}
