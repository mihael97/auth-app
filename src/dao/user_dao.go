package dao

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/mihael97/auth-proxy/src/dto/user"
	"github.com/mihael97/auth-proxy/src/model"
	"gitlab.com/mihael97/Go-utility/src/database"
	"gitlab.com/mihael97/Go-utility/src/util"
	"gitlab.com/mihael97/Go-utility/src/util/mapper"
)

const InsertUser = "INSERT INTO USERS(USERNAME, PASSWORD) VALUES($1, $2) RETURNING ID"
const GetUser = "SELECT * FROM USERS WHERE username = $1 AND is_deleted = false"
const GetUsers = "SELECT * FROM USERS WHERE is_deleted = false"
const DeleteUser = "UPDATE users SET is_deleted = NOT is_deleted WHERE id = $1"

var userDaoImpl *userDao

type userDao struct {
	mapper mapper.DatabaseMapper[model.User]
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

func (r *userDao) mapRow(row *sql.Rows, item *model.User) (err error) {
	err = row.Scan(&item.Id, &item.Username, &item.Password, &item.CreatedOn, &item.IsDeleted)
	return
}

func (r *userDao) GetAllUsers() ([]model.User, error) {
	rows, err := database.GetDatabase().Query(GetUsers)
	if err != nil {
		return nil, err
	}
	return r.mapper.ScanRows(rows, r.mapRow)
}

func (r *userDao) GetUser(username string) (*model.User, error) {
	rows, err := database.GetDatabase().Query(GetUser, username)
	if err != nil {
		return nil, err
	}
	items, err := r.mapper.ScanRows(rows, r.mapRow)
	if err != nil {
		return nil, err
	} else if len(items) == 0 {
		return nil, nil
	}
	return util.GetPointer(items[0]), nil
}

func (d *userDao) CreateUser(request user.CreateUserDto) (*model.User, error) {
	log.Println("Saving user")
	response, err := database.GetDatabase().Query(InsertUser, request.Username, request.Password)
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
			mapper.GetDatabaseMapper(func(item model.User) []any {
				return []any{&item.Id, &item.Username, &item.Password, &item.CreatedOn, &item.IsDeleted}
			}),
		}
	}
	return userDaoImpl
}
