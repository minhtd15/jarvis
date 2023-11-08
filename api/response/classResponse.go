package api_response

import "time"

type CourseInfoResponse struct {
	CourseId      int64     `json:"course_id"`
	CourseName    string    `json:"course_name"`
	MainTeacher   string    `json:"main_teacher"`
	Room          int64     `json:"room"`
	StartDate     string    `json:"start_date"`
	EndDate       string    `json:"end_date"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	StudyDays     string    `json:"study_days"`
	CourseStatus  string    `json:"course_status"`
	TotalSessions int64     `json:"total_sessions"`
}
