package dao

import (
	"github.com/mihael97/auth-proxy/src/dto"
	"github.com/mihael97/auth-proxy/src/model"
)

type UserDao interface {
	CreateUser(request dto.CreateUserDto) (*model.User, error)
	GetUser(username string) (*model.User, error)
}

type CustomerRoleDao interface {
	CreateCustomerRole(id string, roles ...string) error
}
