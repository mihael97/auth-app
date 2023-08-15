package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/mihael97/auth-proxy/src/dao"
	"github.com/mihael97/auth-proxy/src/dto/user"
	"github.com/mihael97/auth-proxy/src/util"
	"gitlab.com/mihael97/Go-utility/src/web/security/jwt"
	"golang.org/x/crypto/bcrypt"
)

const DefaultValidityPerion = uint64(60 * 60 * 1000)

var loginService *loginServiceImpl

type loginServiceImpl struct {
	userRepository dao.UserDao
	maker          jwt.Maker
	validityPeriod time.Duration
}

func (s *loginServiceImpl) Login(request user.LoginUserDto) (*string, error) {
	user, err := s.userRepository.GetUser(request.Username)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, nil
	} else if user.IsDeleted {
		log.Printf("User %s exists but deleted\n", user.Id)
		return nil, nil
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		log.Printf("Error during login for user %s\n%v", request.Username, err)
		if is := errors.As(err, &bcrypt.ErrMismatchedHashAndPassword); is {
			return nil, fmt.Errorf("invalid login")
		}
		return nil, err
	}

	return s.maker.CreateToken(request.Username, s.validityPeriod)
}

func GetLoginService() LoginService {
	if loginService == nil {
		jwtSecret := util.GetConfig().Security.Secret
		if jwtSecret == nil {
			log.Panicln("JWT token not provided")
		}

		validityPeriod := DefaultValidityPerion
		if util.GetConfig().Security.ValidityPeriod != nil {
			validityPeriod = *util.GetConfig().Security.ValidityPeriod
		}

		maker, err := jwt.NewJwtMaker(*jwtSecret)
		if err != nil {
			log.Panicln(err)
		}

		loginService = &loginServiceImpl{
			dao.GetUserDao(),
			maker,
			time.Duration(validityPeriod) * time.Millisecond,
		}
	}
	return loginService
}
