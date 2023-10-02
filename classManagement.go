package education_website

import (
	"context"
	"education-website/entity/course_class"
	"time"
)

type ClassRequest struct {
	ClassName string `json:"class_name"`
}

type ClassInfoResponse struct {
	ClassName      string    `json:"class_name"`
	CourseName     string    `json:"course_name"`
	StudyDate      time.Time `json:"study_date"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Room           string    `json:"room"`
	NumberSessions int       `db:"NUMBER_SESSIONS"`
	Status         int       `json:"status"`
}

type ClassService interface {
	GetClassInformationByClassName(request ClassRequest, ctx context.Context) (ClassInfoResponse, error)
}

type ClassStore interface {
	GetClassInformationByClassNameStore(request ClassRequest, ctx context.Context) (course_class.CourseClassEntity, error)
}
