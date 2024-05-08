package education_website

import (
	"context"
	api_request "education-website/api/request"
	api_response "education-website/api/response"
	"education-website/client/response"
	"education-website/entity/course_class"
	response2 "education-website/rabbitmq/response"
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
	AddNoteService(noteRequest api_request.AddNoteRequest, add []string, delete []string, ctx context.Context) error
	GetTAListService(classId int, ctx context.Context) ([]string, error)
	GetCheckInHistoryByCourseId(courseId string, ctx context.Context) ([]api_response.CheckInHistory, error)
	AddSubClassService(rq api_request.NewSubClassRequest, ctx context.Context) error
	GetSubClassByCourseId(courseId string, ctx context.Context) ([]api_response.SubClassResponse, error)
	DeleteSubClassService(rq api_request.DeleteSubClassRequest, ctx context.Context) error
	GetAllAvailableCourseFeeService() ([]response.CoursesFeeResponse, error)
	GetCourseRevenueByCourseIdService(courseId string, ctx context.Context) (*response.CoursesFeeResponse, error)
	GetCourseByYear(year string, ctx context.Context) error
	UpdateYearlyRevenueAndCourseRevenue(rq response2.YearlyResponse, ctx context.Context) error
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
	AddNoteStore(noteRequest api_request.AddNoteRequest, add []string, delete []string, ctx context.Context) error
	GetTaListInSessionStore(classId int, ctx context.Context) ([]string, error)
	GetCheckInHistoryByCourseIdStore(courseId string, ctx context.Context) ([]course_class.CheckInHistoryEntity, error)
	AddSubClassStore(rq api_request.NewSubClassRequest, ctx context.Context) error
	GetSubClassByCourseIdStore(courseId string, ctx context.Context) ([]course_class.SubClassEntity, error)
	DeleteSubClassStore(rq string, ctx context.Context) error
	GetCourseInfoByCourseId(courseId string, ctx context.Context) (course_class.CourseEntity, error)
	GetTotalStudentByCourseIdStore(courseId string, ctx context.Context) (*int, error)
	GetCourseByYearStore(year string, ctx context.Context) ([]course_class.CourseEntity, error)
	UpdateYearlyRevenueAndCourseRevenue(rq response2.YearlyResponse, ctx context.Context) error
}
