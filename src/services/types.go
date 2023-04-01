package services

import (
	"github.com/mihael97/auth-proxy/src/dto"
	"github.com/mihael97/auth-proxy/src/model"
)

type LoginService interface {
	Login(request dto.LoginUserDto) error
}

type UserService interface {
	CreateUser(request dto.CreateUserDto, username string) (*model.User, error)
}
