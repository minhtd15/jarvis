package education_website

import (
	"context"
	api_request "education-website/api/request"
	api_response "education-website/api/response"
	"education-website/entity/course_class"
	"time"
)

type ClassService interface {
	AddNewClass(request api_request.NewCourseRequest, ctx context.Context) error
	GetCourseInformationByClassName(request api_request.CourseInfoRequest, ctx context.Context) (*api_response.CourseInfoResponse, error)
	GetAllCourses(ctx context.Context) ([]api_response.CourseInfoResponse, error)
}

type ClassStore interface {
	InsertNewCourseStore(entity course_class.CourseEntity, request api_request.NewCourseRequest, schedule []time.Time, ctx context.Context) error
	GetCourseInformationStore(request api_request.CourseInfoRequest, ctx context.Context) (course_class.CourseEntity, error)
	GetSessionsByCourseType(typeCourseCode string, ctx context.Context) (*int, *int, error)
	GetAllCoursesStore(ctx context.Context) ([]course_class.CourseEntity, error)
}
