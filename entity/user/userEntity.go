package user

type UserEntity struct {
	UserId       string `db:"ID"`
	UserName     string `db:"USERNAME"`
	Email        string `db:"EMAIL"`
	Role         string `db:"ROLE"`
	DOB          string `db:"DOB"`
	JobPosition  string `db:"JOB_POSITION"`
	StartingDate string `db:"STARTING_DATE"`
	Password     string `db:"PASSWORD"`
	Gender       string `db:"GENDER"`
	FullName     string `db:"FULL_NAME"`
}
