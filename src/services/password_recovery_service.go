package services

import (
	"github.com/mihael97/auth-proxy/src/dao"
	"github.com/mihael97/auth-proxy/src/dto/passwordRecovery"
)

var passwordRecoveryServiceImpl *passwordRecoveryService

type passwordRecoveryService struct {
	passwordRecoveryDao dao.PasswordRecoveryDao
}

func (p passwordRecoveryService) MarkRecoveryAttemptAsDone(id string) error {
	return p.passwordRecoveryDao.MarkRecoveryAttemptAsDone(id)
}

func (p passwordRecoveryService) RemoveUnusedPasswordRecoveryAttempts(id string) error {
	return p.passwordRecoveryDao.RemoveUnusedPasswordRecoveryAttempts(id)
}

func (p passwordRecoveryService) GetPasswordRecoveryById(id string) (*passwordRecovery.PasswordRecoveryDto, error) {
	passwordRecoveryAttempt, err := p.passwordRecoveryDao.GetPasswordRecoveryById(id)
	if err != nil {
		return nil, err
	}
	return &passwordRecovery.PasswordRecoveryDto{
		Id:             passwordRecoveryAttempt.Id,
		UserId:         passwordRecoveryAttempt.UserId,
		CreationDate:   passwordRecoveryAttempt.CreationDate,
		ExpirationDate: passwordRecoveryAttempt.ExpirationDate,
		ExecutionDate:  passwordRecoveryAttempt.ExecutionDate,
		Unused:         passwordRecoveryAttempt.Unused,
	}, nil
}

func (p passwordRecoveryService) CreatePasswordRecoveryAttempt(username string) (*string, error) {
	return p.passwordRecoveryDao.CreatePasswordRecoveryAttempt(username)
}

func (p passwordRecoveryService) IsPasswordRecoveryActive(username string) (bool, error) {
	return p.passwordRecoveryDao.IsPasswordRecoveryActive(username)
}

func GetPasswordRecoveryService() PasswordRecoveryService {
	if passwordRecoveryServiceImpl == nil {
		passwordRecoveryServiceImpl = &passwordRecoveryService{
			passwordRecoveryDao: dao.GetPasswordRecoveryDao(),
		}
	}
	return passwordRecoveryServiceImpl
}
