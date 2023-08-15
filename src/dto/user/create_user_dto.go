package user

type CreateUserDto struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
	Email    *string  `json:"email"`
}
