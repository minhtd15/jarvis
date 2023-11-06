package api_response

type SalaryStoreResponse struct {
	UserName    string            `json:"user_name"`
	FullName    string            `json:"full_name"`
	Gender      string            `json:"gender"`
	JobPosition string            `json:"job_position"`
	Salary      SalaryInformation `json:"salary"`
}

type SalaryInformation struct {
	PayrollId  int64   `json:"payroll_id"`
	CourseType string  `json:"course_type"`
	WorkDays   int64   `json:"work_days"`
	PriceEach  float64 `json:"price_each"`
	Amount     float64 `json:"amount"`
}

type SalaryAPIResponse struct {
	UserId      string              `json:"user_id"`
	UserName    string              `json:"user_name"`
	FullName    string              `json:"full_name"`
	Gender      string              `json:"gender"`
	JobPosition string              `json:"job_position"`
	Salary      []SalaryInformation `json:"salary"`
}
