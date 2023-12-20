package request

import "time"

type NewCourseRequest struct {
	CourseType  string   `json:"class_type"` // PS, SR, MT
	StartTime   string   `json:"start_time"` // gio bat dau hoc
	EndTime     string   `json:"end_time"`   // gio ket thuc ca hoc
	StartDate   string   `json:"start_date"`
	MainTeacher string   `json:"teacher"`
	Room        int      `json:"room"`
	StudyDays   []string `json:"study_date"` // ngay hoc (vd: t2,4,6)
	Location    string   `json:"location"`   // co so hoc: hqv, vvd, dcv, hdt
}

type CourseInfoRequest struct {
	CourseId string `json:"course_id"`
}

type DeleteClassInfo struct {
	ClassId string `json:"class_id"`
}

type ModifyCourseInformation struct {
	CourseId int    `json:"course_id"`
	Teacher  string `json:"teacher"`
	Room     int    `json:"room"`
}

type AddNoteRequest struct {
	ClassId   string    `json:"class_id"`
	Date      string    `json:"date"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Room      int       `json:"room"`
	Note      string    `json:"note"`
}
