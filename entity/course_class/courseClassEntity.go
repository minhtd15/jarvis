package course_class

import "database/sql"

type CourseClassEntity struct {
	ClassId        string       `db:"CLASS_ID"`
	ClassName      string       `db:"CLASS_NAME"`
	CourseId       string       `db:"COURSE_ID"`
	CourseName     string       `db:"COURSE_NAME"`
	StudyDate      sql.NullTime `db:"STUDY_DATE"`
	StartDate      sql.NullTime `db:"START_DATE"`
	EndDate        sql.NullTime `db:"END_DATE"`
	Room           int          `db:"ROOM"`
	NumberSessions int          `db:"NUMBER_SESSIONS"`
}
