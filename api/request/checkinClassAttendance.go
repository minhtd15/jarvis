package request

import "time"

type ChecInAttendanceClassRequest struct {
	UserId      string    `json:"user_id"`
	ClassId     string    `json:"class_id"`
	CourseType  string    `json:"course_type"`
	CheckInUser string    `json:"check_in_user"`
	CheckInTime time.Time `json:"check_in_time"`
}
