package salary

import "time"

type AttendanceHistoryEntity struct {
	AttendanceId string    `db:"ATTENDANCE_ID"` // lưu trữ giá trị duy nhất của từng record mỗi khi checkin
	UserId       string    `db:"USER_ID"`       // lưu giá trị của người thực hiện checkin
	ClassName    string    `db:"CLASS_NAME"`    // lưu giá trị của lớp checkin
	CourseType   string    `db:"COURSE_TYPE"`   // lưu giá trị của loại công việc: trông thi, thư viện, inclass,...
	CheckInTime  time.Time `db:"CHECKIN_TIME"`
	Status       string    `db:"STATUS"`
}
