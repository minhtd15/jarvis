package request

type CheckInAttendanceWorkerRequest struct {
	CourseId     string `json:"course_id"`
	ClassId      string `json:"class_id"`
	CourseTypeId string `json:"course_type_id"`
}

type CheckInWorkerHistory struct {
	CourseId string `json:"courseId"`
}
