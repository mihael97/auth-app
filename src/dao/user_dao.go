package dao

import (
	"log"
	"time"

	"github.com/mihael97/auth-proxy/src/dto"
	"github.com/mihael97/auth-proxy/src/model"
	"gitlab.com/mihael97/Go-utility/src/database"
	"gitlab.com/mihael97/Go-utility/src/util/mapper"
)

const InsertUser = "INSERT INTO USERS(USERNAME, PASSWORD) VALUES($1, $2) RETURNING ID"
const GetUser = "SELECT * FROM USERS WHERE username = $1"

var userDaoImpl *userDao

type userDao struct {
	mapper mapper.DatabaseMapper[model.User]
}

func (r *userDao) GetUser(username string) (*model.User, error) {
	rows, err := database.GetDatabase().Query(GetUser, username)
	if err != nil {
		return nil, err
	}
	return r.mapper.MapItem(rows)
}

func (d *userDao) CreateUser(request dto.CreateUserDto) (*model.User, error) {
	log.Println("Saving user")
	response, err := database.GetDatabase().Query(InsertUser, request.Username, request.Password)
	if err != nil {
		return nil, err
	}
	var id string
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
			mapper.GetDatabaseMapper(func(item model.User) []any {
				return []any{&item.Id, &item.Username, &item.Password, &item.CreatedOn, &item.IsDeleted}
			}),
		}
	}
	return userDaoImpl
}
