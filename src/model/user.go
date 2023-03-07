package model

import "time"

type User struct {
	Id        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedOn time.Time `json:"createdOn"`
	IsDeleted bool      `json:"isDeleted"`
	Roles     []string  `json:"roles"`
}
