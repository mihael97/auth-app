package services

import (
	"fmt"
	"log"
	"strings"

	"github.com/mihael97/auth-proxy/src/dao"
	"github.com/mihael97/auth-proxy/src/dto/user"
	"github.com/mihael97/auth-proxy/src/model"
	"gitlab.com/mihael97/Go-utility/src/util"
)

var userService *userServiceImpl

type userServiceImpl struct {
	userRepository  dao.UserDao
	customerRoleDao dao.CustomerRoleDao
}

func (s *userServiceImpl) CreateUser(request user.CreateUserDto, username string) (*user.UserDto, error) {
	if len(request.Roles) == 0 {
		request.Roles = []string{string(model.USER)}
	} else if util.Contains(model.ADMIN.String(), request.Roles...) {
		currentUser, err := s.userRepository.GetUser(username)
		if err != nil {
			return nil, err
		}
		if !util.Contains(model.ADMIN.String(), currentUser.Roles...) {
			return nil, fmt.Errorf("cannot create ADMIN user with logged user. Please check privilages")
		}

		// append user role
		if !util.Contains(model.USER.String(), request.Roles...) {
			request.Roles = append(request.Roles, model.USER.String())
		}
	}
	createdUser, err := s.userRepository.CreateUser(request)
	if err != nil {
		return nil, err
	}
	log.Printf("Created user %s\n", createdUser.Id)

	err = s.customerRoleDao.CreateCustomerRole(createdUser.Id, request.Roles...)
	if err != nil {
		return nil, err
	}
	log.Printf("Created roles (%s) for %s\n", strings.Join(request.Roles, ","), username)

	createdUser.Roles = request.Roles
	return &user.UserDto{
		Id:        createdUser.Id,
		Username:  createdUser.Username,
		CreatedOn: createdUser.CreatedOn,
		IsDeleted: createdUser.IsDeleted,
		Roles:     createdUser.Roles,
	}, nil
}

func GetUserService() UserService {
	if userService == nil {
		userService = &userServiceImpl{
			dao.GetUserDao(),
			dao.GetCustomerRoleDao(),
		}
	}
	return userService
}
