package model

import "time"

type PasswordRecovery struct {
	Id             string     `json:"id"`
	CreationDate   time.Time  `json:"creationDate"`
	ExpirationDate time.Time  `json:"expirationDate"`
	ExecutionDate  *time.Time `json:"executionDate"`
	UserId         string     `json:"userId"`
	Unused         bool       `json:"unused"`
}
