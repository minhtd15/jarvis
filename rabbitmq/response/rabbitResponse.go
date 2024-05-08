package response

type CourseFeeResponse struct {
	CourseID int     `json:"courseID"`
	Revenue  float64 `json:"revenue"`
}

// YearlyResponse struct tương đương với lớp YearlyResponse trong Java
type YearlyResponse struct {
	CourseFeeResponses []CourseFeeResponse `json:"courseFeeResponses"`
	TotalYearlyRevenue float64             `json:"totalYearlyRevenue"`
	Year               string              `json:"year"`
}
