package panther

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type userStore struct {
	db *sqlx.DB
}

func NewUserStore(db *sqlx.DB) UserStore {
	return &userStore{
		db: db,
	}
}

func (u userStore) GetByUserName(userName string) (UserEntity, error) {
	ctx := context.Background()
	logrus.Infof("Get user information by UserName")

	entity := UserEntity{}
	sqlQuery := "SELECT * FROM USER WHERE USER_NAME = ?"

	// execute sql query
	err := u.db.QueryRowxContext(ctx, sqlQuery, userName).Scan(&entity.Id, &entity.UserName, &entity.Password, &entity.PhoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.WithError(err).Errorf("Cannot find user with user name: %s", userName)
			return entity, err
		}
		logrus.WithError(err).Errorf("Cannot get info from database for user: %s", userName)
		return entity, err
	}
	return entity, nil
}
