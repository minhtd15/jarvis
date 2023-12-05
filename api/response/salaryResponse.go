package api_response

type SalaryStoreResponse struct {
	UserName    string            `json:"user_name"`
	FullName    string            `json:"full_name"`
	Gender      string            `json:"gender"`
	JobPosition string            `json:"job_position"`
	Salary      SalaryInformation `json:"salary"`
}

type SalaryInformation struct {
	PayrollId int64   `json:"payroll_id"`
	WorkDays  int64   `json:"work_days"`
	Amount    float64 `json:"amount"`
}

type SalaryConfig struct {
	PayrollId   int64   `json:"payroll_id"`
	PayrollRate float64 `json:"payroll_rate"`
	TypePayroll string  `json:"course_type"`
}

type SalaryAPIResponse struct {
	UserId       string              `json:"user_id"`
	UserName     string              `json:"user_name"`
	FullName     string              `json:"full_name"`
	Gender       string              `json:"gender"`
	JobPosition  string              `json:"job_position"`
	SalaryConfig []SalaryConfig      `json:"salary_config"`
	Salary       []SalaryInformation `json:"salary"`
}
