package user

import "time"

type UserDto struct {
	Id        string    `json:"id"`
	Username  string    `json:"username"`
	CreatedOn time.Time `json:"createdOn"`
	IsDeleted bool      `json:"isDeleted"`
	Roles     []string  `json:"roles"`
}
