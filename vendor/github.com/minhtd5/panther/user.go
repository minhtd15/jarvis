package panther

import (
	_ "database/sql"
)

type UserStore interface {
	GetByUserName(userName string) (UserEntity, error)
}

type UserService interface {
	GetByUserName(userName string) (UserEntity, error)
}

type UserEntity struct {
	Id          int64  `db:"ID"` // id = service_name + uuid + date
	PhoneNumber string `db:"PHONE_NUMBER"`
	UserName    string `db:"USER_NAME"`
	Password    string `db:"Password"`
}
