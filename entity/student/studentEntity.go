package student

import "time"

type EntityStudent struct {
	Id      string    `db:"STUDENT_ID"`
	Name    string    `db:"NAME"`
	DOB     time.Time `db:"DOB"`
	Email   string    `db:"EMAIL"`
	PhoneNo string    `db:"PHONE_NUMBER"`
}

type StudentAttendanceEntity struct {
	StudentAttendanceId int64  `db:"STUDENT_ATTENDANCE_ID"`
	StudentId           int64  `db:"STUDENT_ID"`
	ClassId             int64  `db:"CLASS_ID"`
	Status              string `db:"STATUS"`
}
