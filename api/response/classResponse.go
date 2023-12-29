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

type ClassResponse struct {
	ClassId   string   `json:"class_id"`
	StartTime string   `json:"start_time"`
	EndTime   string   `json:"end_time"`
	Date      string   `json:"date"`
	Room      string   `json:"room"`
	TypeClass string   `json:"type_class"`
	Note      string   `json:"note"`
	TaList    []string `json:"assistant"`
}

type SubClassResponse struct {
	ClassId   string `json:"class_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Date      string `json:"date"`
	Room      string `json:"room"`
	TaId      string `json:"ta_id"`
	Note      string `json:"note"`
}
