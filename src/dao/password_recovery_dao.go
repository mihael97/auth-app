package dao

import (
	"fmt"
	"github.com/mihael97/auth-proxy/src/model"
	"gitlab.com/mihael97/Go-utility/src/database"
	"gitlab.com/mihael97/Go-utility/src/util"
	"gitlab.com/mihael97/Go-utility/src/util/mapper"
	"log"
	"time"
)

const GetPasswordRecoveryAttemptById = "SELECT * FROM auth.PASSWORD_RECOVERY WHERE id = $1"
const InsertPasswordRecovery = "INSERT INTO auth.PASSWORD_RECOVERY(USER_ID, EXPIRATION_DATE) VALUES ((SELECT u.ID FROM auth.users u WHERE u.username=$1), $2) RETURNING ID"
const GetNumberOfActivePasswordRecoveries = "SELECT COUNT(*) as recoveries FROM auth.PASSWORD_RECOVERY WHERE USER_ID = (SELECT u.ID FROM auth.users u WHERE u.username=$1) AND EXPIRATION_DATE IS NULL"
const RemoveUnusedPasswordRecoveryAttempts = "UPDATE auth.PASSWORD_RECOVERY SET unused=true WHERE user_id=$1 AND execution_time IS NULL AND expiration_date < current_timestamp"
const ExecutePasswordRecovery = "UPDATE auth.PASSWORD_RECOVERY SET unused=false, execution_time=current_timestamp WHERE id=$1"

var passwordRecoveryDaoImpl *passwordRecoveryDao

type passwordRecoveryDao struct {
	GmtLocation    *time.Location
	databaseMapper mapper.DatabaseMapper[model.PasswordRecovery]
}

func (p passwordRecoveryDao) MarkRecoveryAttemptAsDone(id string) error {
	_, err := database.GetDatabase().Query(ExecutePasswordRecovery, id)
	return err
}

func (p passwordRecoveryDao) RemoveUnusedPasswordRecoveryAttempts(id string) error {
	_, err := database.GetDatabase().Query(RemoveUnusedPasswordRecoveryAttempts, id)
	return err
}

func (p passwordRecoveryDao) GetPasswordRecoveryById(id string) (*model.PasswordRecovery, error) {
	rows, err := database.GetDatabase().Query(GetPasswordRecoveryAttemptById, id)
	if err != nil {
		return nil, err
	}
	items, err := p.databaseMapper.MapItems(rows)
	if err != nil {
		return nil, err
	} else if len(items) == 0 {
		return nil, nil
	}
	return &items[0], nil
}

func (p passwordRecoveryDao) CreatePasswordRecoveryAttempt(username string) (*string, error) {
	expirationDate := time.Now().In(p.GmtLocation).Add(24 * time.Hour)
	result, err := database.GetDatabase().Query(InsertPasswordRecovery, username, expirationDate)
	if err != nil {
		return nil, err
	}
	var id string
	if !result.Next() {
		return nil, fmt.Errorf("no rows found")
	}
	if err := result.Scan(&id); err != nil {
		return nil, err
	}
	return util.GetPointer(id), nil
}

func (p passwordRecoveryDao) IsPasswordRecoveryActive(username string) (bool, error) {
	rows, err := database.GetDatabase().Query(GetNumberOfActivePasswordRecoveries, username)
	if err != nil {
		return false, err
	}
	rows.Next()
	var count string
	err = rows.Scan(&count)
	if err != nil {
		return false, err
	}
	return count != "0", nil
}

func GetPasswordRecoveryDao() PasswordRecoveryDao {
	if passwordRecoveryDaoImpl == nil {
		gmtLocation, err := time.LoadLocation("Europe/London")
		if err != nil {
			log.Panic(err)
		}
		passwordRecoveryDaoImpl = &passwordRecoveryDao{
			GmtLocation: gmtLocation,
			databaseMapper: mapper.GetDatabaseMapper(func(rows mapper.SqlRowsData) model.PasswordRecovery {
				creationDate := (*rows.GetData("creation_time")).(time.Time)
				expirationDate := (*rows.GetData("expiration_date")).(time.Time)
				var executionDate *time.Time
				executionDateValue := rows.GetData("execution_date")
				if executionDateValue != nil {
					executionDate = (*executionDateValue).(*time.Time)
				}

				return model.PasswordRecovery{
					Id:             rows.GetString("id"),
					CreationDate:   creationDate,
					ExpirationDate: expirationDate,
					ExecutionDate:  executionDate,
					UserId:         rows.GetString("user_id"),
					Unused:         rows.GetBool("unused"),
				}
			}),
		}
	}
	return passwordRecoveryDaoImpl
}
