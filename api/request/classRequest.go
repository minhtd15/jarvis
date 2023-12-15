package request

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
