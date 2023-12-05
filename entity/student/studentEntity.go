package student

import "time"

type EntityStudent struct {
	Id      string    `db:"STUDENT_ID"`
	Name    string    `db:"NAME"`
	DOB     time.Time `db:"DOB"`
	Email   string    `db:"EMAIL"`
	PhoneNo string    `db:"PHONE_NUMBER"`
}
