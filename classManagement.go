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
	GetCourseType(ctx context.Context) (map[int]string, error)
	GetFromToSchedule(fromDate string, toDate string, userId string, courseType map[int]string, ctx context.Context) ([]api_response.FromToScheduleResponse, error)
	DeleteCourseByCourseId(courseId string, ctx context.Context) error
	DeleteClassByClassId(rq api_request.DeleteClassInfo, ctx context.Context) error
	GetAllSessionsByCourseIdService(courseId string, ctx context.Context) ([]api_response.ClassResponse, error)
	FixCourseInformationService(rq api_request.ModifyCourseInformation, ctx context.Context) error
	AddNoteService(noteRequest api_request.AddNoteRequest, ctx context.Context) error
}

type ClassStore interface {
	InsertNewCourseStore(entity course_class.CourseEntity, request api_request.NewCourseRequest, schedule []time.Time, ctx context.Context) error
	GetCourseInformationStore(request api_request.CourseInfoRequest, ctx context.Context) (course_class.CourseEntity, error)
	GetSessionsByCourseType(typeCourseCode string, ctx context.Context) (*int, *int, error)
	GetAllCoursesStore(ctx context.Context) ([]course_class.CourseEntity, error)
	GetAllCourseType(ctx context.Context) ([]course_class.CourseTypeEntity, error)
	GetClassFromToDateStore(fromDate string, toDate string, userId string, ctx context.Context) ([]course_class.FromToScheduleEntity, error)
	DeleteCourseById(courseId string, ctx context.Context) error
	DeleteClassByIdStore(classId string, ctx context.Context) error
	GetAllSessionsByCourseIdStore(courseId string, ctx context.Context) ([]course_class.ClassEntity, error)
	FixCourseInformationStore(rq api_request.ModifyCourseInformation, ctx context.Context) error
	AddNoteStore(noteRequest api_request.AddNoteRequest, ctx context.Context) error
}
