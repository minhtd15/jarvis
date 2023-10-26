package course_class

import "database/sql"

type CourseEntity struct {
	CourseId      int64
	CourseTypeId  int64
	MainTeacher   string
	Room          int64
	StartDate     sql.NullString
	EndDate       sql.NullString
	StudyDays     string
	CourseName    string
	TotalSessions int64
}
