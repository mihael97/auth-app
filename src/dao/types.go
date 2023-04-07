package dao

import (
	"github.com/mihael97/auth-proxy/src/dto/user"
	"github.com/mihael97/auth-proxy/src/model"
)

type UserDao interface {
	CreateUser(request user.CreateUserDto) (*model.User, error)
	GetUser(username string) (*model.User, error)
	GetAllUsers() ([]model.User, error)
}

type CustomerRoleDao interface {
	CreateCustomerRole(id string, roles ...string) error
	GetUserRoles(id string) ([]string, error)
}
