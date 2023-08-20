package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/mihael97/auth-proxy/src/dto/passwordRecovery"
	"log"
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
	userService  services.UserService
	loginService services.LoginService
	headerName   string
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

func (u *userControllerImpl) getUserInfo(ctx *gin.Context) {
	currentUsername := u.parseUsername(ctx)
	if currentUsername == nil {
		log.Println("username is null")
		web.ParseToJson(pointerUtil.CreateValidationException(map[string][]string{"username": {"username is empty"}}), ctx, http.StatusBadRequest)
		return
	}
	userDto, err := u.userService.GetUser(*currentUsername)
	if err != nil {
		web.WriteError(err, ctx)
		return
	}
	web.ParseToJson(userDto, ctx, http.StatusOK)
}

func (u *userControllerImpl) getUsers(ctx *gin.Context) {
	users, err := u.userService.GetUsers()
	if err != nil {
		web.WriteError(err, ctx)
		return
	}
	web.ParseToJson(users, ctx, http.StatusOK)
}

func (u *userControllerImpl) deleteUser(ctx *gin.Context) {
	var request user.DeleteUserDto
	if err := json.NewDecoder(ctx.Request.Body).Decode(&request); err != nil {
		web.WriteError(err, ctx)
		return
	}
	err := u.userService.DeleteUser(request.Id, *u.parseUsername(ctx))
	if err != nil {
		web.WriteError(err, ctx)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (u *userControllerImpl) parseUsername(ctx *gin.Context) *string {
	username, err := jwt.GetUserNameFromToken(ctx, *util.GetConfig().Security.Secret)
	if err != nil {
		fmt.Println("Error during parsing username: ", err)
		return nil
	}
	return pointerUtil.GetPointer(username)
}

func (u *userControllerImpl) sendRecoveryEmail(context *gin.Context) {
	var request user.SendPasswordRecoveryDto

	if err := json.NewDecoder(context.Request.Body).Decode(&request); err != nil {
		web.WriteError(err, context)
		return
	}

	err := u.userService.SendRecoveryEmail(request)
	if err != nil {
		web.WriteError(err, context)
		return
	}
	context.Status(http.StatusNoContent)
}

func (u *userControllerImpl) passwordRecovery(context *gin.Context) {
	var request passwordRecovery.PasswordRecoveryRequest
	if err := json.NewDecoder(context.Request.Body).Decode(&request); err != nil {
		web.WriteError(err, context)
		return
	}

	username, err := u.userService.ChangePassword(request)
	if err != nil {
		web.WriteError(err, context)
		return
	}

	u.loginUser(*username, request.NewPassword, context)
}

func (u *userControllerImpl) loginUser(username, password string, ctx *gin.Context) {
	loginRequest := user.LoginUserDto{
		Username: username,
		Password: password,
	}
	token, err := u.loginService.Login(loginRequest)
	if err != nil {
		web.WriteError(err, ctx)
		return
	}

	ctx.Writer.Header().Add(u.headerName, *token)
	ctx.Status(http.StatusNoContent)
}

func GetUserController() *gin.Engine {
	if userController == nil {
		headerName := util.GetConfig().Security.HeaderName
		if headerName == nil {
			headerName = pointerUtil.GetPointer("Authorization")
		}
		userController = &userControllerImpl{
			services.GetUserService(),
			services.GetLoginService(),
			*headerName,
		}
	}

	engine := gin.New()
	group := engine.Group("/api/users")
	group.POST("/", userController.createUser)
	group.GET("/me", userController.getUserInfo)
	group.GET("/", userController.getUsers)
	group.DELETE("/", userController.deleteUser)
	group.POST("/recovery", userController.passwordRecovery)
	group.POST("/recovery/email", userController.sendRecoveryEmail)
	return engine
}
