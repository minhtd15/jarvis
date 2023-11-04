package api_response

type FromToScheduleResponse struct {
	CourseId   int64  `json:"course_id"`
	CourseCode string `json:"course_code"`
	CourseName string `json:"course_name"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	Date       string `json:"date"`
}
