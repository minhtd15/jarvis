package user

type UserEntity struct {
	UserId       string `db:"user_id"`
	UserName     string `db:"user_name"`
	Email        string `db:"email"`
	Role         string `db:"Role"`
	DOB          string `db:"dob"`
	JobPosition  string `db:"job_position"`
	StartingDate string `db:"starting_date"`
	Password     string `db:"password"`
}
