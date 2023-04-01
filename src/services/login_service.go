package services

import "github.com/mihael97/auth-proxy/src/dto/user"

var loginService *loginServiceImpl

type loginServiceImpl struct {
}

func (*loginServiceImpl) Login(request user.LoginUserDto) error {
	panic("unimplemented")
}

func GetLoginService() LoginService {
	if loginService == nil {
		loginService = &loginServiceImpl{}
	}
	return loginService
}
