package salary

import "time"

type SalaryEntity struct {
	SalaryId   int        `db:"SALARY_ID"`
	EmployeeId string     `db:"EMPLOYEE_ID"`
	Month      time.Month `db:"MONTH"`
	TypeJobId  int        `db:"TYPE_JOB_ID"`
	WorkDays   int        `db:"WORD_DAYS"`
}
