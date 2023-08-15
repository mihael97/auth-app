package dao

import (
	"github.com/mihael97/auth-proxy/src/dto/passwordRecovery"
	"github.com/mihael97/auth-proxy/src/dto/user"
	"github.com/mihael97/auth-proxy/src/model"
)

type UserDao interface {
	CreateUser(request user.CreateUserDto) (*model.User, error)
	GetUser(username string) (*model.User, error)
	GetUserById(id string) (*model.User, error)
	GetAllUsers() ([]model.User, error)
	DeleteUser(id string) error
	ChangePassword(id string, request passwordRecovery.PasswordRecoveryRequest) error
}

type CustomerRoleDao interface {
	CreateCustomerRole(id string, roles ...string) error
	GetUserRoles(id string) ([]string, error)
}

type PasswordRecoveryDao interface {
	CreatePasswordRecoveryAttempt(username string) (*string, error)
	IsPasswordRecoveryActive(username string) (bool, error)
	GetPasswordRecoveryById(id string) (*model.PasswordRecovery, error)
	RemoveUnusedPasswordRecoveryAttempts(id string) error
	MarkRecoveryAttemptAsDone(id string) error
}
