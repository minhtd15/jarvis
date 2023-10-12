package request

type ModifySalaryConfRequest struct {
	UserId        string                `json:"user_id"`
	UserName      string                `json:"user_name"`
	NewSalaryList []SalaryConfiguration `json:"new_salary_list"`
}

type SalaryConfiguration struct {
	PayrollId     string `json:"payroll_id"`
	TypePayroll   string `json:"type_payroll"`
	PayrollAmount string `json:"payroll_amount"`
}
