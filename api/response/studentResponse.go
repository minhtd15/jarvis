package api_response

import "time"

type StudentResponse struct {
	StudentId   string    `json:"student_id"`
	StudentName string    `json:"student_name"`
	Dob         time.Time `json:"dob"`
}
