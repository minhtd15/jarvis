package education_website

import (
	"context"
	api_request "education-website/api/request"
	api_response "education-website/api/response"
	"education-website/entity/course_class"
)

type ClassService interface {
	//AddNewClass(request api_request.NewClassRequest, ctx context.Context) error
	GetCourseInformationByClassName(request api_request.CourseInfoRequest, ctx context.Context) (*api_response.CourseInfoResponse, error)
}

type ClassStore interface {
	InsertNewClassStore(request api_request.NewClassRequest, ctx context.Context) error
	GetCourseInformationStore(request api_request.CourseInfoRequest, ctx context.Context) (course_class.CourseEntity, error)
}
