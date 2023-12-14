package request

type StudentAttendanceRequest struct {
	StudentId int64  `json:"student_id"`
	ClassId   int64  `json:"class_id"`
	Status    string `json:"status"`
}
