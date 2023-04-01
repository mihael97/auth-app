package services

import "github.com/mihael97/auth-proxy/src/dto"

var loginService *loginServiceImpl

type loginServiceImpl struct {
}

func (*loginServiceImpl) Login(request dto.LoginUserDto) error {
	panic("unimplemented")
}

func GetLoginService() LoginService {
	if loginService == nil {
		loginService = &loginServiceImpl{}
	}
	return loginService
}
