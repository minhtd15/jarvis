package api_response

import "time"

type StudentResponse struct {
	StudentId   string    `json:"student_id"`
	StudentName string    `json:"student_name"`
	Dob         time.Time `json:"dob"`
}

type StudentAttendanceScheduleResponse struct {
	StudentId int64            `json:"student_id"`
	Name      string           `json:"name"`
	Dob       string           `json:"dob"`
	CheckIn   []CheckInStudent `json:"checkIn"`
}

type CheckInStudent struct {
	ClassId int64  `json:"class_id"`
	Date    string `json:"date"`
	Status  string `json:"status"`
}
