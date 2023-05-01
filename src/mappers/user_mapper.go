package mappers

import (
	"log"

	"github.com/mihael97/auth-proxy/src/dao"
	"github.com/mihael97/auth-proxy/src/dto/user"
	"github.com/mihael97/auth-proxy/src/model"
	"gitlab.com/mihael97/Go-utility/src/util/mapper"
)

type UserMapper struct {
	mapper.GenericMapper[model.User, user.UserDto]
	customerRoleDao dao.CustomerRoleDao
}

func (g *UserMapper) MapItem(row model.User) (dto *user.UserDto) {
	dto = &user.UserDto{
		Id:        row.Id,
		Username:  row.Username,
		CreatedOn: row.CreatedOn,
		IsDeleted: row.IsDeleted,
	}
	var err error
	dto.Roles, err = g.customerRoleDao.GetUserRoles(row.Id)
	if err != nil {
		log.Println(err)
		return nil
	}
	return dto
}

func (g *UserMapper) MapItems(rows []model.User) []user.UserDto {
	users := make([]user.UserDto, len(rows))

	for i := range rows {
		mappedValue := g.MapItem(rows[i])
		if mappedValue == nil {
			return nil
		}
		users[i] = *mappedValue
	}

	return users
}

func GetUserMapper() UserMapper {
	return UserMapper{
		customerRoleDao: dao.GetCustomerRoleDao(),
	}
}
