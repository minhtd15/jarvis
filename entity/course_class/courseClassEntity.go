package course_class

import "database/sql"

type CourseEntity struct {
	CourseId      int            `db:"COURSE_ID"`
	CourseTypeId  int            `db:"COURSE_TYPE_ID"`
	MainTeacher   string         `db:"MAIN_TEACHER"`
	Room          int            `db:"ROOM"`
	StartDate     sql.NullString `db:"START_DATE"`
	EndDate       sql.NullString `db:"END_DATE"`
	StartTime     sql.NullString `db:"START_TIME"`
	EndTime       sql.NullString `db:"END_TIME"`
	StudyDays     string         `db:"STUDY_DAYS"`
	CourseName    string         `db:"COURSE_NAME"`
	TotalSessions int64          `db:"TOTAL_SESSIONS"`
	Location      sql.NullString `db:"LOCATION"`
}

type CheckInHistoryEntity struct {
	UserId      string `db:"USER_ID"`
	ClassId     int    `db:"CLASS_ID"`
	CheckInTime string `db:"CHECKIN_TIME"`
	Status      string `db:"STATUS"`
}
