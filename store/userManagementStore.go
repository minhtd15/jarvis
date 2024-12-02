package store

import (
	"context"
	"database/sql"
	batman "education-website"
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
	log.Infof("Retrieving user information for UserName: %s", userName)

	entity := batman.UserResponse{}
	sqlQuery := "SELECT * FROM USER WHERE USERNAME = ? OR EMAIL = ? OR USER_ID = ?"
	//var tmp sql.NullString∂

	// execute sql query
	err := u.db.QueryRowxContext(ctx, sqlQuery, userName, email, userId).Scan(
		&entity.UserId,
		&entity.UserName,
		&entity.Email,
		&entity.Role,
		&entity.DOB,
		&entity.StartDate,
		&entity.JobPosition,
		&entity.Password,
		&entity.FullName,
		&entity.Gender,
		//&tmp,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Infof("No user found with UserName: %s, Email: %s, or UserID: %d", userName, email, userId)
			return entity, nil
		}
		log.WithError(err).Errorf("Failed to get user info from database for UserName: %s", userName)
		return entity, err
	}

	return entity, nil

}

func (u *userManagementStore) InsertNewUserStore(newUser user.UserEntity, ctx context.Context) error {
	log.Infof("insert new user to database after validated register information")

	// Start a new transaction
	tx, err := u.db.Begin()
	if err != nil {
		log.WithError(err).Errorf("Failed to begin transaction")
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			// Something went wrong, rollback the transaction
			tx.Rollback()
			panic(p) // Re-throw the panic after rollback
		} else if err != nil {
			// Error occurred, rollback the transaction
			tx.Rollback()
		}
	}()

	sqlQuery := "INSERT INTO USER VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	stmt, err := tx.Prepare(sqlQuery)
	if err != nil {
		log.WithError(err).Errorf("Failed to prepare SQL statement")
		return err
	}
	defer stmt.Close()

	tmp := 0
	// Execute the prepared statement
	result, err := stmt.Exec(newUser.UserId, newUser.UserName, newUser.Email, newUser.Role, newUser.DOB, newUser.StartingDate, newUser.JobPosition, newUser.Password, newUser.FullName, newUser.Gender, tmp)
	if err != nil {
		log.WithError(err).Errorf("Failed to insert user into the database")
		return err
	}

	// Get the last insert ID
	_, err = result.LastInsertId()
	if err != nil {
		log.WithError(err).Errorf("Failed to get last insert ID")
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.WithError(err).Errorf("Failed to commit transaction")
		return err
	}

	return nil
}

func (u *userManagementStore) UpdateNewPassword(newPassword []byte, userName string) error {
	log.Infof("Start to update new password")

	// Begin a transaction
	tx, err := u.db.Begin()
	if err != nil {
		log.WithError(err).Errorf("Failed to begin transaction")
		return err
	}

	defer func() {
		// Rollback the transaction if there is an error or return is not nil
		if r := recover(); r != nil || err != nil {
			log.WithError(err).Errorf("Rolling back transaction")
			tx.Rollback()
			return
		}
		// Commit the transaction if there is no error
		err := tx.Commit()
		if err != nil {
			log.WithError(err).Errorf("Failed to commit transaction")
		}
	}()

	sqlQuery := "UPDATE USER SET PASSWORD = ? WHERE USERNAME = ?"

	stmt, err := tx.Prepare(sqlQuery)
	if err != nil {
		log.WithError(err).Errorf("Failed to prepare SQL statement")
		return err
	}
	defer stmt.Close()

	// Execute the prepared statement within the transaction
	_, err = stmt.Exec(newPassword, userName)
	if err != nil {
		log.WithError(err).Errorf("Failed to update user password in the database")
		return err
	}

	// Return nil if the update is successful
	return nil
}

func (u *userManagementStore) GetSalaryReportStore(userName string, month string, year string, ctx context.Context) ([]salary.SalaryEntity, error) {
	log.Infof("Get salary information from database")

	var entities []salary.SalaryEntity

	var sqlQuery string
	var rows *sqlx.Rows
	sqlQuery = "select u.USER_ID, u.USERNAME, u.FULLNAME, u.GENDER, u.JOB_POSITION, s.PAYROLL_ID, s.TYPE_PAYROLL, s.TOTAL_WORK_DATES, s.PAYROLL_RATE, s.SALARY " +
		"from SALARY s join USER u ON s.USER_ID = u.USER_ID " +
		"WHERE s.MONTH = ? and s.YEAR = ?"
	args := []interface{}{month, year}

	if userName != "" && userName != "undefined" {
		sqlQuery += " AND u.FULLNAME LIKE CONCAT('%', ?, '%')"
		args = append(args, userName)
	}

	rows, err := u.db.QueryxContext(ctx, sqlQuery, args...)
	if err != nil {
		log.WithError(err).Errorf("Cannot get info from the database for user: %s", userName)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var entity salary.SalaryEntity
		if err := rows.Scan(&entity.UserId, &entity.UserName, &entity.FullName, &entity.Gender, &entity.JobPosition, &entity.PayrollId, &entity.TypeWork, &entity.TotalWorkDates, &entity.PayrollPerSessions, &entity.TotalSalary); err != nil {
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
