package api_response

type CourseResponse struct {
	CourseId   string `json:"course_id"`
	CourseName string `json:"course_name"`
	Room       string `json:"room"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	StudyDays  string `json:"study_days"`
	Location   string `json:"location"`
}

type CheckInHistory struct {
	UserId      string `json:"userId"`
	ClassId     int    `json:"classId"`
	CheckInTime string `json:"checkInTime"`
	Status      string `json:"status"`
}

type CourseFeeResponse struct {
	CourseId      string  `json:"courseId"`
	FeePerStudent float64 `json:"feePerStudent"`
	TotalFee      float64 `json:"totalFee"`
}
