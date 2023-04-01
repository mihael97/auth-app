package services

import "github.com/mihael97/auth-proxy/src/dto/user"

type LoginService interface {
	Login(request user.LoginUserDto) error
}

type UserService interface {
	CreateUser(request user.CreateUserDto, username string) (*user.UserDto, error)
}
