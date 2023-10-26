package request

type NewClassRequest struct {
	ClassType   string `json:"class_type"`
	StartTime   string `json:"start_time"` // gio bat dau hoc
	EndTime     string `json:"end_time"`   // gio ket thuc
	OpeningDate string `json:"opening_date"`
	Teacher     string `json:"teacher"`
	Room        string `json:"room"`
	StudyDate   string `json:"study_date"` // ngay hoc (vd: t2,4,6)
	TypeClass   string `json:"type_class"` // bo tro hay hoc chinh
}

type CourseInfoRequest struct {
	CourseId string `json:"course_id"`
}
