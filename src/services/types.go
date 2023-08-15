package services

import (
	"github.com/mihael97/auth-proxy/src/dto/passwordRecovery"
	"github.com/mihael97/auth-proxy/src/dto/user"
)

type LoginService interface {
	Login(request user.LoginUserDto) (*string, error)
}

type UserService interface {
	CreateUser(request user.CreateUserDto, username string) (*user.UserDto, error)
	GetUser(username string) (*user.UserDto, error)
	GetUsers() ([]user.UserDto, error)
	DeleteUser(id, username string) error
	SendRecoveryEmail(request user.SendPasswordRecoveryDto) error
	ChangePassword(request passwordRecovery.PasswordRecoveryRequest) (*string, error)
}

type PasswordRecoveryService interface {
	CreatePasswordRecoveryAttempt(username string) (*string, error)
	IsPasswordRecoveryActive(username string) (bool, error)
	GetPasswordRecoveryById(id string) (*passwordRecovery.PasswordRecoveryDto, error)
	RemoveUnusedPasswordRecoveryAttempts(id string) error
	MarkRecoveryAttemptAsDone(id string) error
}
