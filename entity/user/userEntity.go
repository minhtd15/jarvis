package user

type UserEntity struct {
	UserId       string `db:"ID"`
	UserName     string `db:"USERNAME"`
	Email        string `db:"email"`
	Role         string `db:"Role"`
	DOB          string `db:"dob"`
	JobPosition  string `db:"job_position"`
	StartingDate string `db:"starting_date"`
	Password     string `db:"password"`
}
