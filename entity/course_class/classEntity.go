package course_class

import (
	"database/sql"
)

type ClassEntity struct {
	CourseId  string         `db:"COURSE_ID"`
	ClassId   string         `db:"CLASS_ID"`
	StartTime string         `db:"START_TIME"`
	EndTime   string         `db:"END_TIME"`
	Date      sql.NullString `db:"DATE"`
	Room      string         `db:"ROOM"`
	TypeClass string         `db:"TYPE_CLASS"`
	Note      sql.NullString `db:"NOTE"`
}

type SubClassEntity struct {
	ClassId   string         `db:"CLASS_ID"`
	StartTime string         `db:"START_TIME"`
	EndTime   string         `db:"END_TIME"`
	Date      sql.NullString `db:"DATE"`
	Room      string         `db:"ROOM"`
	TaId      string         `db:"USER_ID"`
	Note      sql.NullString `db:"NOTE"`
}
