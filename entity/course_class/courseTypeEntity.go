package course_class

type CourseTypeEntity struct {
	CourseTypeId  int    `db:"COURSE_TYPE_ID"`
	Name          string `db:"NAME"`
	CourseCode    string `db:"CODE"`
	TotalSessions string `db:"TOTAL_SESSIONS"`
}
