package mappers

import (
	"log"

	"github.com/mihael97/auth-proxy/src/dao"
	"github.com/mihael97/auth-proxy/src/dto/user"
	"github.com/mihael97/auth-proxy/src/model"
	"gitlab.com/mihael97/Go-utility/src/util/mapper"
)

type UserMapper struct {
	mapper.Mapper[model.User, user.UserDto]
	customerRoleDao dao.CustomerRoleDao
}

func GetUserMapper() mapper.Mapper[model.User, user.UserDto] {
	return mapper.GetDefaultMapper(func(item model.User) *user.UserDto {
		mappedItem := &user.UserDto{
			Id:        item.Id,
			Username:  item.Username,
			CreatedOn: item.CreatedOn,
			IsDeleted: item.IsDeleted,
			Email:     item.Email,
		}
		var err error
		item.Roles, err = dao.GetCustomerRoleDao().GetUserRoles(item.Id)
		if err != nil {
			log.Println(err)
			return nil
		}
		return mappedItem
	})
}
