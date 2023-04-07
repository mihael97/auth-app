package services

import (
	"fmt"
	"log"
	"strings"

	"github.com/mihael97/auth-proxy/src/dao"
	"github.com/mihael97/auth-proxy/src/dto/user"
	"github.com/mihael97/auth-proxy/src/model"
	"gitlab.com/mihael97/Go-utility/src/util"
	"gitlab.com/mihael97/Go-utility/src/util/mapper"
)

var userService *userServiceImpl

type userServiceImpl struct {
	userRepository  dao.UserDao
	customerRoleDao dao.CustomerRoleDao
	dtoMapper       mapper.Mapper[model.User, user.UserDto]
}

func (s *userServiceImpl) GetUser(username string) (*user.UserDto, error) {
	fetchedUser, err := s.userRepository.GetUser(username)
	if err != nil {
		return nil, err
	} else if fetchedUser == nil {
		return nil, nil
	}
	fetchedUser.Roles, err = s.customerRoleDao.GetUserRoles(fetchedUser.Id)
	if err != nil {
		return nil, err
	}
	return s.dtoMapper.MapItem(*fetchedUser), nil
}

func (s *userServiceImpl) CreateUser(request user.CreateUserDto, username string) (*user.UserDto, error) {
	//check if already exists
	existingUser, err := s.userRepository.GetUser(request.Username)
	if err != nil {
		return nil, err
	} else if existingUser != nil {
		return nil, fmt.Errorf("username already exists")
	}

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
			mapper.GetGenericMapper(func(item model.User) user.UserDto {
				return user.UserDto{
					Id:        item.Id,
					Username:  item.Username,
					CreatedOn: item.CreatedOn,
					IsDeleted: item.IsDeleted,
					Roles:     item.Roles,
				}
			}),
		}
	}
	return userService
}
