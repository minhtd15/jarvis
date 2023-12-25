package request

type DeleteStudentRequest struct {
	StudentId string `json:"student_id"`
	CourseId  string `json:"course_id"`
}
