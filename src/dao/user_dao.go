package dao

import (
	"fmt"
	"github.com/mihael97/auth-proxy/src/dto/passwordRecovery"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"

	"github.com/mihael97/auth-proxy/src/dto/user"
	"github.com/mihael97/auth-proxy/src/model"
	"gitlab.com/mihael97/Go-utility/src/database"
	"gitlab.com/mihael97/Go-utility/src/util"
	"gitlab.com/mihael97/Go-utility/src/util/mapper"
)

const InsertUser = "INSERT INTO auth.users(USERNAME, PASSWORD) VALUES($1, $2) RETURNING ID"
const InsertUserWithEmail = "INSERT INTO auth.users(USERNAME, PASSWORD, EMAIL) VALUES($1, $2, $3) RETURNING ID"
const GetUser = "SELECT * FROM auth.users WHERE username = $1 AND is_deleted = false"
const GetUserById = "SELECT * FROM auth.users WHERE id = $1 AND is_deleted = false"
const GetUsers = "SELECT * FROM auth.users WHERE is_deleted = false"
const DeleteUser = "UPDATE auth.users SET is_deleted = NOT is_deleted WHERE id = $1"
const UpdatePassword = "UPDATE auth.users SET password = $2 WHERE id = $1"

var userDaoImpl *userDao

type userDao struct {
	mapper mapper.DatabaseMapper[model.User]
}

func (r *userDao) GetUserById(id string) (*model.User, error) {
	rows, err := database.GetDatabase().Query(GetUserById, id)
	if err != nil {
		return nil, err
	}
	items, err := r.mapper.MapItems(rows)
	if err != nil {
		return nil, err
	} else if len(items) == 0 {
		return nil, nil
	}
	return util.GetPointer(items[0]), nil
}

func (*userDao) ChangePassword(id string, request passwordRecovery.PasswordRecoveryRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), 14)
	if err != nil {
		return err
	}
	result, err := database.GetDatabase().Exec(UpdatePassword, id, hashedPassword)

	if err != nil {
		return err
	} else if rowsAffected, _ := result.RowsAffected(); rowsAffected != 1 {
		return fmt.Errorf("wrong number of affected rows")
	}

	return nil
}

func (*userDao) DeleteUser(id string) error {
	result, err := database.GetDatabase().Exec(DeleteUser, id)
	if err != nil {
		return err
	} else if rowsAffected, _ := result.RowsAffected(); rowsAffected != 1 {
		return fmt.Errorf("wrong number of rows affected")
	}
	return nil
}

func (r *userDao) GetAllUsers() ([]model.User, error) {
	rows, err := database.GetDatabase().Query(GetUsers)
	if err != nil {
		return nil, err
	}
	return r.mapper.MapItems(rows)
}

func (r *userDao) GetUser(username string) (*model.User, error) {
	rows, err := database.GetDatabase().Query(GetUser, username)
	if err != nil {
		return nil, err
	}
	items, err := r.mapper.MapItems(rows)
	if err != nil {
		return nil, err
	} else if len(items) == 0 {
		return nil, nil
	}
	return util.GetPointer(items[0]), nil
}

func (d *userDao) CreateUser(request user.CreateUserDto) (*model.User, error) {
	log.Println("Saving user")

	query := InsertUser
	args := []string{request.Username, request.Password}
	if request.Email != nil {
		args = append(args, *request.Email)
		query = InsertUserWithEmail
	}

	response, err := database.GetDatabase().Query(query, args)
	if err != nil {
		return nil, err
	}
	var id string
	if !response.Next() {
		return nil, fmt.Errorf("invalid database response")
	}
	err = response.Scan(&id)
	if err != nil {
		return nil, err
	}
	log.Printf("Created user %s\n", id)
	return &model.User{
		Id:        id,
		Username:  request.Username,
		Password:  request.Password,
		CreatedOn: time.Now(),
		IsDeleted: false,
	}, nil
}

func GetUserDao() UserDao {
	if userDaoImpl == nil {
		userDaoImpl = &userDao{
			mapper.GetDatabaseMapper(func(rows mapper.SqlRowsData) model.User {
				var email *string
				if rows.HasValue("email") {
					email = util.GetPointer(rows.GetString("email"))
				}
				createdOn := (*rows.GetData("created_on")).(time.Time)

				return model.User{
					Id:        rows.GetString("id"),
					Username:  rows.GetString("username"),
					Password:  rows.GetString("password"),
					CreatedOn: createdOn,
					IsDeleted: rows.GetBool("is_deleted"),
					Roles:     nil,
					Email:     email,
				}
			}),
		}
	}
	return userDaoImpl
}
