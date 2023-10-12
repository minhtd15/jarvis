package salary

type SalaryEntity struct {
	UserName           string  `db:"USERNAME"`
	FullName           string  `db:"FULLNAME"`
	Gender             string  `db:"GENDER"`
	JobPosition        string  `db:"JOB_POSITION"`
	TypeWork           string  `db:"TYPE_PAYROLL"`
	TotalWorkDates     int64   `db:"TOTAL_WORK_DATES"`
	PayrollPerSessions float64 `db:"INCOME"`
	TotalSalary        float64 `db:"SALARY"`
}
