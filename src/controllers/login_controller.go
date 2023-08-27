package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mihael97/auth-proxy/src/dto/user"
	"github.com/mihael97/auth-proxy/src/security"
	"github.com/mihael97/auth-proxy/src/services"
	exceptionUtil "gitlab.com/mihael97/Go-utility/src/util"
	"gitlab.com/mihael97/Go-utility/src/web"
	"net/http"
)

var loginControllerImpl *loginController

type loginController struct {
	loginService services.LoginService
}

// @Summary Login user
// @Tags login
// @Param request body user.LoginUserDto true "Login request"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /login [post]
func (c *loginController) loginUser(ctx *gin.Context) {
	var request user.LoginUserDto
	if err := json.NewDecoder(ctx.Request.Body).Decode(&request); err != nil {
		web.WriteError(err, ctx)
		return
	}
	token, err := c.loginService.Login(request)
	if err != nil {
		web.WriteError(err, ctx)
		return
	} else if token == nil {
		exception := exceptionUtil.NewException(
			fmt.Sprintf("user %s not found", request.Username),
		)
		web.ParseToJson(exception, ctx, http.StatusUnauthorized)
		return
	}
	ctx.Writer.Header().Add(security.AuthorizationHeader, *token)
	appendExpiresAt(ctx)
	ctx.Request.Header.Del(security.AuthorizationHeader)
	ctx.Status(http.StatusNoContent)
}

func GetLoginController() *gin.Engine {
	if loginControllerImpl == nil {
		loginControllerImpl = &loginController{loginService: services.GetLoginService()}
	}
	engine := gin.New()

	group := engine.Group("/api/login")
	group.POST("/", loginControllerImpl.loginUser)
	return engine
}
