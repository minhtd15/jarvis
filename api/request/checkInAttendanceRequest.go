package request

type CheckInAttendanceWorkerRequest struct {
	UserId       string `json:"user_id"`
	CourseId     string `json:"course_id"`
	ClassId      string `json:"class_id"`
	CourseTypeId string `json:"course_type_id"`
}

type CheckInWorkerHistory struct {
	CourseId string `json:"courseId"`
}
