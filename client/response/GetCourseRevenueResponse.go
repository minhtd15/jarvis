package response

type GetCourseRevenueByCourseIdResponse struct {
	CourseId string `yaml:"courseId"`
}

type CoursesFeeResponse struct {
	CourseId      string  `json:"courseId"`
	CourseTypeId  int     `json:"courseTypeId"`
	FeePerStudent float64 `json:"feePerStudent"`
	TotalStudent  int     `json:"totalStudent"`
	TotalFee      float64 `json:"totalFee"`
}

type PaymentStatusByCourseIdResponse struct {
	CourseManagerId int     `json:"courseManagerId"`
	StudentId       int     `json:"studentId"`
	Remain          float64 `json:"remain"`
	PaymentStatus   string  `json:"paymentStatus"`
}

type YearReportRevenueResponse struct {
	Message string `json:"message"`
}
