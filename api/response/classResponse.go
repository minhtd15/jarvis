package api_response

import "time"

type CourseInfoResponse struct {
	CourseId      int64     `json:"course_id"`
	CourseName    string    `json:"course_name"`
	MainTeacher   string    `json:"main_teacher"`
	Room          int64     `json:"room"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	StudyDays     string    `json:"study_days"`
	CourseStatus  string    `json:"course_status"`
	TotalSessions int64     `json:"total_sessions"`
}
