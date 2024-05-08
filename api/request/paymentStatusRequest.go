package request

import "education-website/entity/course_class"

type PaymentStatusRequest struct {
	CourseTypeId         int                                `json:"courseTypeId"`
	CourseManagerRequest []course_class.CourseManagerEntity `json:"courseManagerRequest"`
}
