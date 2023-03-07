package controllers

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/mihael97/auth-proxy/src/dto"
	"gitlab.com/mihael97/Go-utility/src/web"
	"gitlab.com/mihael97/Go-utility/src/web/routes"
)

var loginControllerImpl *loginController

type loginController struct {
}

func (*loginController) GetBasePath() string {
	return "/login"
}

func (c *loginController) GetRoutes() map[routes.Route]func(ctx *gin.Context) {
	return map[routes.Route]func(ctx *gin.Context){
		routes.CreateRoute("/", web.POST, false): c.loginUser,
	}
}

func (c *loginController) loginUser(ctx *gin.Context) {
	var request dto.CreateUserDto
	if err := json.NewDecoder(ctx.Request.Body).Decode(&request); err != nil {
		web.WriteError(err, ctx)
		return
	}
}

func GetLoginController() routes.RoutesController {
	if loginControllerImpl == nil {
		loginControllerImpl = &loginController{}
	}
	return loginControllerImpl
}
