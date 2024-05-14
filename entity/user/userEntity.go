package user

import "database/sql"

type UserEntity struct {
	UserId        string        `db:"USER_ID"`
	UserName      string        `db:"USERNAME"`
	Email         string        `db:"EMAIL"`
	Role          string        `db:"ROLE"`
	DOB           string        `db:"DOB"`
	JobPosition   string        `db:"JOB_POSITION"`
	StartingDate  string        `db:"STARTINGDATE"`
	Password      string        `db:"PASSWORD"`
	Gender        string        `db:"GENDER"`
	FullName      string        `db:"FULLNAME"`
	ResetPassword sql.NullInt64 `db:"RESET_PASSWORD"`
}
