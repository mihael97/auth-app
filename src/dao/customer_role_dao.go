package dao

import (
	"log"

	"gitlab.com/mihael97/Go-utility/src/database"
)

const InsertCustomerRole = "INSERT INTO customer_roles(ROLE_NAME, USER_ID) VALUES ($1, $2)"
const GetUserRoles = "SELECT ROLE_NAME FROM customer_roles WHERE USER_ID = $1"

var customerRoleDao *customerRoleDaoImpl

type customerRoleDaoImpl struct {
}

func (*customerRoleDaoImpl) GetUserRoles(id string) ([]string, error) {
	rows, err := database.GetDatabase().Query(GetUserRoles, id)
	if err != nil {
		return nil, err
	}
	var roles []string
	var roleName string

	for rows.Next() {
		if err := rows.Scan(&roleName); err != nil {
			return nil, err
		}
		roles = append(roles, roleName)
	}

	return roles, nil
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
