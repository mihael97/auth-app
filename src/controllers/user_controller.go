package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mihael97/auth-proxy/src/dto/user"
	"github.com/mihael97/auth-proxy/src/services"
	"github.com/mihael97/auth-proxy/src/util"
	pointerUtil "gitlab.com/mihael97/Go-utility/src/util"
	"gitlab.com/mihael97/Go-utility/src/web"
	"gitlab.com/mihael97/Go-utility/src/web/security/jwt"
)

var userController *userControllerImpl

type userControllerImpl struct {
	userService services.UserService
	headerName  string
}

func (u *userControllerImpl) createUser(ctx *gin.Context) {
	var request user.CreateUserDto
	if err := json.NewDecoder(ctx.Request.Body).Decode(&request); err != nil {
		web.WriteError(err, ctx)
		return
	}

	createdUser, err := u.userService.CreateUser(request, *u.parseUsername(ctx))
	if err != nil {
		web.WriteError(err, ctx)
		return
	}
	web.ParseToJson(createdUser, ctx, http.StatusCreated)
}

func (u *userControllerImpl) parseUsername(ctx *gin.Context) *string {
	username, err := jwt.GetUserNameFromToken(ctx, *util.GetConfig().Security.Secret)
	if err != nil {
		fmt.Println("Error during parsing username: ", err)
		return nil
	}
	return pointerUtil.GetPointer(username)
}

func GetUserController() *gin.Engine {
	if userController == nil {
		headerName := util.GetConfig().Security.HeaderName
		if headerName == nil {
			headerName = pointerUtil.GetPointer("Authorization")
		}
		userController = &userControllerImpl{
			services.GetUserService(),
			*headerName,
		}
	}

	engine := gin.New()
	group := engine.Group("/api/users")
	group.POST("/", userController.createUser)
	return engine
}