package dao

import (
	"log"

	"gitlab.com/mihael97/Go-utility/src/database"
)

const InsertCustomerRole = "INSERT INTO customer_roles(ROLE_NAME, USER_ID) VALUES ($1, $2)"

var customerRoleDao *customerRoleDaoImpl

type customerRoleDaoImpl struct {
}

func (*customerRoleDaoImpl) CreateCustomerRole(id string, roles ...string) error {
	tx, err := database.GetDatabase().Begin()
	if err != nil {
		return err
	}
	for _, role := range roles {
		_, err = tx.Exec(InsertCustomerRole, role, id)
		if err != nil {
			log.Printf("Error during adding role %s for %s\n", role, id)
			return err
		}
	}
	return tx.Commit()
}

func GetCustomerRoleDao() CustomerRoleDao {
	if customerRoleDao == nil {
		customerRoleDao = &customerRoleDaoImpl{}
	}
	return customerRoleDao
}
