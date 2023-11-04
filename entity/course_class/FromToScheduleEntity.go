package course_class

type FromToScheduleEntity struct {
	CourseId     int64  `db:"COURSE_ID"`
	CourseTypeId int64  `db:"COURSE_TYPE_ID"`
	StartTime    string `db:"START_TIME"`
	EndTime      string `db:"END_TIME"`
	Date         string `db:"DATE"`
}
