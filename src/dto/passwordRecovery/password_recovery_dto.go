package passwordRecovery

import (
	"time"
)

type PasswordRecoveryDto struct {
	Id             string     `json:"id"`
	UserId         string     `json:"userId"`
	CreationDate   time.Time  `json:"creationDate"`
	ExpirationDate time.Time  `json:"expirationDate"`
	ExecutionDate  *time.Time `json:"executionDate"`
	Unused         bool       `json:"unused"`
}
