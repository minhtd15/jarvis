package course_class

import "time"

type ClassEntity struct {
	ClassId   string    `db:"CLASS_ID"`
	StartTime time.Time `db:"START_TIME"`
	EndTime   time.Time `db:"END_TIME"`
}
