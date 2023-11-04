package request

type ModifySalaryConfRequest struct {
	UserId        string                `json:"user_id"`
	NewSalaryList []SalaryConfiguration `json:"salary"`
}

type SalaryConfiguration struct {
	PayrollId     int     `json:"payroll_id"`
	TypePayroll   string  `json:"type_payroll"`
	PayrollAmount float64 `json:"payroll_amount"`
}
