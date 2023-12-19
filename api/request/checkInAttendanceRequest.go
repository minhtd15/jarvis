package request

type CheckInAttendanceWorkerRequest struct {
	CourseId     string `json:"course_id"`
	CourseName   string `json:"class_id"`
	CourseTypeId string `json:"course_type_id"`
}
