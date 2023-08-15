package services

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/mihael97/auth-proxy/src/dto/passwordRecovery"
	config "github.com/mihael97/auth-proxy/src/util"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/go-mail/mail"
	"github.com/mihael97/auth-proxy/src/dao"
	"github.com/mihael97/auth-proxy/src/dto/user"
	"github.com/mihael97/auth-proxy/src/mappers"
	"github.com/mihael97/auth-proxy/src/model"
	"gitlab.com/mihael97/Go-utility/src/util"
)

var userService *userServiceImpl

type userServiceImpl struct {
	userRepository          dao.UserDao
	passwordRecoveryService PasswordRecoveryService
	customerRoleDao         dao.CustomerRoleDao
	dtoMapper               mappers.UserMapper
}

func (s *userServiceImpl) ChangePassword(request passwordRecovery.PasswordRecoveryRequest) (*string, error) {
	recoveryAttempt, err := s.passwordRecoveryService.GetPasswordRecoveryById(request.AttemptId)
	if err != nil {
		return nil, err
	} else if recoveryAttempt == nil {
		return nil, fmt.Errorf("recovery attempt doesn't exist")
	}

	if err = validateRecoveryAttempt(recoveryAttempt); err != nil {
		return nil, err
	}

	fetchedUser, err := s.userRepository.GetUserById(recoveryAttempt.UserId)
	if err != nil {
		return nil, err
	}

	err = checkIfPasswordMatch(fetchedUser.Password, request.NewPassword)
	if err != nil {
		return nil, err
	}

	//change passwords
	err = s.passwordRecoveryService.MarkRecoveryAttemptAsDone(request.AttemptId)
	if err != nil {
		return nil, err
	}
	err = s.userRepository.ChangePassword(fetchedUser.Id, request)
	if err != nil {
		return nil, err
	}

	return util.GetPointer(fetchedUser.Username), nil
}

func checkIfPasswordMatch(currentPassword, newPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(currentPassword), []byte(newPassword))
	if err != nil {
		if errors.As(err, &bcrypt.ErrMismatchedHashAndPassword) {
			return nil
		}
		return err
	}
	return fmt.Errorf("password is the same")
}

func validateRecoveryAttempt(recoveryAttempt *passwordRecovery.PasswordRecoveryDto) error {
	if recoveryAttempt.ExpirationDate.Before(time.Now().In(time.UTC)) {
		return fmt.Errorf("recovery attempt not active")
	} else if recoveryAttempt.ExecutionDate != nil {
		return fmt.Errorf("recovery attempt already used")
	}
	return nil
}

func (s *userServiceImpl) SendRecoveryEmail(request user.SendPasswordRecoveryDto) error {
	if len(request.Username) == 0 {
		return fmt.Errorf("username is empty")
	}

	isActive, err := s.passwordRecoveryService.IsPasswordRecoveryActive(request.Username)
	if err != nil {
		return err
	} else if isActive {
		return fmt.Errorf("user %s has active password recovery", request.Username)
	}

	fetchedUser, err := s.userRepository.GetUser(request.Username)
	if err != nil {
		return err
	} else if fetchedUser == nil {
		return fmt.Errorf("user doesn't exist")
	} else if fetchedUser.Email == nil {
		return fmt.Errorf("email not set for %s", request.Username)
	}

	err = s.passwordRecoveryService.RemoveUnusedPasswordRecoveryAttempts(fetchedUser.Id)
	if err != nil {
		return err
	}

	attemptId, err := s.passwordRecoveryService.CreatePasswordRecoveryAttempt(request.Username)
	if err != nil {
		return err
	}
	link := fmt.Sprintf("%s/recovery/%s", *config.GetConfig().PasswordRecovery.Url, *attemptId)
	err = s.sendEmail(fetchedUser, link)
	return err
}

func (s *userServiceImpl) DeleteUser(id string, username string) error {
	currentUser, err := s.GetUser(username)
	if err != nil {
		return err
	} else if currentUser == nil {
		return fmt.Errorf("user not found")
	} else if !util.Contains("ADMIN", currentUser.Roles...) {
		return fmt.Errorf("unsufficient role")
	}
	return s.userRepository.DeleteUser(id)
}

func (s *userServiceImpl) GetUsers() ([]user.UserDto, error) {
	users, err := s.userRepository.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return s.dtoMapper.MapItems(users), nil
}

func (s *userServiceImpl) GetUser(username string) (*user.UserDto, error) {
	fetchedUser, err := s.userRepository.GetUser(username)
	if err != nil {
		return nil, err
	} else if fetchedUser == nil {
		return nil, nil
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

func (s *userServiceImpl) sendEmail(u *model.User, link string) error {
	log.Println("Sending email to", *u.Email, "for password recovery for user", u.Username)

	emailConfig := config.GetConfig().EmailConfig

	from := *emailConfig.From
	password := *emailConfig.ServerConfig.Secret
	to := []string{
		*u.Email,
	}

	m := mail.NewMessage()
	m.SetHeaders(map[string][]string{
		"From":    {from},
		"To":      to,
		"Subject": {"Password recovery"},
	})

	emailHost := *emailConfig.ServerConfig.Host
	emailPort, _ := strconv.ParseInt(*emailConfig.ServerConfig.Port, 10, 64)

	bodyData := make(map[string]interface{}, 0)
	bodyData["User"] = u
	bodyData["URL"] = link
	body, err := s.renderBody(bodyData)
	if err != nil {
		return err
	}

	m.SetBody("text/html", *body)

	d := mail.NewDialer(emailHost, int(emailPort), from, password)
	err = d.DialAndSend(m)
	return err
}

func (s *userServiceImpl) renderBody(u map[string]interface{}) (*string, error) {
	t, err := template.ParseFiles("templates/recovery_email_template.html")
	if err != nil {
		return nil, err
	}
	writer := bytes.Buffer{}
	err = t.Execute(&writer, u)
	if err != nil {
		return nil, err
	}
	return util.GetPointer(writer.String()), nil
}

func GetUserService() UserService {
	if userService == nil {
		userService = &userServiceImpl{
			dao.GetUserDao(),
			GetPasswordRecoveryService(),
			dao.GetCustomerRoleDao(),
			mappers.GetUserMapper(),
		}
	}
	return userService
}
