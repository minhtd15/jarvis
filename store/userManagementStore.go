package store

import (
	"context"
	"database/sql"
	batman "education-website"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type userManagementStore struct {
	db *sqlx.DB
}

func (u *userManagementStore) GetByUserNameStore(userName string, ctx context.Context) (batman.UserResponse, error) {
	log.Infof("Get user information by UserName")

	entity := batman.UserResponse{}
	sqlQuery := "SELECT * FROM USER WHERE USER_NAME = ?"

	// execute sql query
	err := u.db.QueryRowxContext(ctx, sqlQuery, userName).Scan(&entity.UserId, &entity.UserName, &entity.DOB, &entity.JobPosition)
	if err != nil {
		if err == sql.ErrNoRows {
			log.WithError(err).Errorf("Cannot find user with user name: %s", userName)
			return entity, err
		}
		log.WithError(err).Errorf("Cannot get info from database for user: %s", userName)
		return entity, err
	}
	return entity, nil
}
