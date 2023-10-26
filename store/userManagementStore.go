package store

import (
	"context"
	"database/sql"
	batman "education-website"
	api_request "education-website/api/request"
	"education-website/entity/salary"
	"education-website/entity/user"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type userManagementStore struct {
	db *sqlx.DB
}

type UserManagementStoreCfg struct {
	Db *sqlx.DB
}

func NewUserManagementStore(userManagementStoreCfg UserManagementStoreCfg) *userManagementStore {
	return &userManagementStore{
		db: userManagementStoreCfg.Db,
	}
}

func (u *userManagementStore) GetByUserNameStore(userName string, email string, userId string, ctx context.Context) (batman.UserResponse, error) {
	log.Infof("Get user information by UserName")

	entity := batman.UserResponse{}
	sqlQuery := "SELECT * FROM USER WHERE USERNAME = ? OR EMAIL = ? OR USER_ID = ?"

	// execute sql query
	err := u.db.QueryRowxContext(ctx, sqlQuery, userName, email, userId).Scan(&entity.UserId, &entity.UserName, &entity.Email, &entity.Role, &entity.DOB, &entity.StartDate, &entity.JobPosition, &entity.Password, &entity.FullName, &entity.Gender)
	if err != nil {
		if err == sql.ErrNoRows {
			//log.WithError(err).Errorf("Cannot find user with user name: %s", userName)
			return entity, nil
		}
		log.WithError(err).Errorf("Cannot get info from database for user: %s", userName)
		return entity, err
	}
	return entity, nil
}

func (u *userManagementStore) InsertNewUserStore(newUser user.UserEntity, ctx context.Context) error {
	log.Infof("insert new user to database after validated register information")

	sqlQuery := "INSERT INTO USER VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	stmt, err := u.db.Prepare(sqlQuery)
	if err != nil {
		log.WithError(err).Errorf("Failed to prepare SQL statement")
		return err
	}
	defer stmt.Close()

	// Execute the prepared statement
	_, err = stmt.Exec(newUser.UserId, newUser.UserName, newUser.Email, newUser.Role, newUser.DOB, newUser.StartingDate, newUser.JobPosition, newUser.Password, newUser.FullName, newUser.Gender)
	if err != nil {
		log.WithError(err).Errorf("Failed to insert user into the database")
		return err
	}
	return nil
}

func (u *userManagementStore) UpdateNewPassword(newPassword []byte, userName string) error {
	log.Infof("Start to update new password")

	sqlQuery := "UPDATE USER SET PASSWORD = ? WHERE USERNAME = ?"

	stmt, err := u.db.Prepare(sqlQuery)
	if err != nil {
		log.WithError(err).Errorf("Failed to prepare SQL statement")
		return err
	}
	defer stmt.Close()

	// Execute the prepared statement
	_, err = stmt.Exec(newPassword, userName)
	if err != nil {
		log.WithError(err).Errorf("Failed to insert user into the database")
		return err
	}
	return nil
}

func (u *userManagementStore) GetSalaryReportStore(userName string, month string, year string, ctx context.Context) ([]salary.SalaryEntity, error) {
	log.Infof("Get salary information from database")

	var entities []salary.SalaryEntity

	var sqlQuery string
	if userName == "" {
		sqlQuery = "select u.USERNAME, u.FULLNAME, u.GENDER, u.JOB_POSITION, s.PAYROLL_ID, s.TYPE_PAYROLL, s.TOTAL_WORK_DATES, s.PAYROLL_RATE, s.SALARY from SALARY s  join USER u ON s.USER_ID = u.USER_ID WHERE s.MONTH = ?  and s.YEAR = ?"
	} else {
		sqlQuery = "select u.USERNAME, u.FULLNAME, u.GENDER, u.JOB_POSITION, s.PAYROLL_ID, s.TYPE_PAYROLL, s.TOTAL_WORK_DATES, s.PAYROLL_RATE, s.SALARY from SALARY s  join USER u ON s.USER_ID = u.USER_ID WHERE s.MONTH = ?  and s.YEAR = ? AND u.FULLNAME LIKE CONCAT('%', ?, '%');"
	}
	rows, err := u.db.QueryxContext(ctx, sqlQuery, month, year, userName)
	if err != nil {
		log.WithError(err).Errorf("Cannot get info from database for user: %s", userName)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entity salary.SalaryEntity
		if err := rows.Scan(&entity.UserName, &entity.FullName, &entity.Gender, &entity.JobPosition, &entity.TypeWork, &entity.TotalWorkDates, &entity.PayrollPerSessions, &entity.TotalSalary); err != nil {
			log.WithError(err).Errorf("Error scanning row: %s", err.Error())
			return nil, err
		}
		entities = append(entities, entity)
	}

	if err := rows.Err(); err != nil {
		log.WithError(err).Errorf("Error iterating rows: %s", err.Error())
		return nil, err
	}

	return entities, nil
}

func (u *userManagementStore) ModifySalaryConfigurationStore(userId string, userSalaryInfo []api_request.SalaryConfiguration, ctx context.Context) error {
	log.Infof("modify salary information for user %s", userId)

	sqlQuery := "UPDATE EMPLOYEE_RATE SET PAYROLL_RATE = ? WHERE USER_ID = ? AND PAYROLL_ID = ?"

	for _, entities := range userSalaryInfo {
		_, err := u.db.ExecContext(ctx, sqlQuery, entities.PayrollAmount, userId, entities.PayrollId)
		if err != nil {
			log.WithError(err).Errorf("Failed to update salary configuration for user: %s", userId)
			return err
		}
	}

	log.Infof("Salary configuration updated successfully for user: %s", userId)
	return nil
}
