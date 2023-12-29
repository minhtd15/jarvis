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

type ModifyCourseInformation struct {
	CourseId int    `json:"course_id"`
	Teacher  string `json:"teacher"`
	Room     int    `json:"room"`
}

type AddNoteRequest struct {
	ClassId   string   `json:"class_id"`
	Date      string   `json:"date"`
	StartTime string   `json:"startTime"`
	EndTime   string   `json:"endTime"`
	Room      int      `json:"room"`
	TaList    []string `json:"assistant"`
	Note      string   `json:"note"`
	Check     bool     `json:"check"`
}

type NewSubClassRequest struct {
	CourseId  int    `json:"course_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Date      string `json:"date"`
	Room      int    `json:"room"`
	TaId      string `json:"ta_id"`
}

type DeleteSubClassRequest struct {
	ClassId string `json:"class_id"`
}
