package model

import "time"

type RoleName string

func (role RoleName) String() string {
	return string(role)
}

const USER RoleName = "USER"
const ADMIN RoleName = "ADMIN"

type User struct {
	Id        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedOn time.Time `json:"createdOn"`
	IsDeleted bool      `json:"isDeleted"`
	Roles     []string  `json:"roles"`
	Email     *string   `json:"email"`
}
