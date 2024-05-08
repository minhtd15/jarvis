package course_class

type CourseManagerEntity struct {
	CourseManagerId int `db:"COURSE_MANAGER_ID" json:"courseManagerId"`
	CourseId        int `db:"COURSE_ID" json:"courseId"`
	StudentId       int `db:"STUDENT_ID" json:"studentId"`
}
