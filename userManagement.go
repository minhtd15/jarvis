package education_website

import (
	"context"
	api_request "education-website/api/request"
	api_response "education-website/api/response"
	"education-website/client/response"
	"education-website/entity/course_class"
	"education-website/entity/salary"
	"education-website/entity/student"
	"education-website/entity/user"
	"mime/multipart"
)

type UserRequest struct {
	Id       string `json:"id"` // id = service_name + uuid + date
	Email    string `json:"email"`
	Dob      string `json:"dob"`
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type UserResponse struct {
	UserId      string `json:"user_id"`
	UserName    string `json:"user_name"`
	DOB         string `json:"dob"`
	Email       string `json:"email"`
	JobPosition string `json:"job_position"`
	Role        string `json:"role"`
	StartDate   string `json:"start_date"`
	Password    string `json:"password"`
	FullName    string `json:"full_name"`
	Gender      string `json:"gender"`
}

type UserService interface {
	GetByUserName(userName string, email string, userId string, ctx context.Context) (*UserResponse, error)
	GetUserNamePassword(userLoginInfo api_request.LoginRequest, ctx context.Context) (*UserResponse, error)
	InsertNewUser(userRegisterInfo api_request.RegisterRequest, ctx context.Context) (string, error)
	ChangePassword(changePasswordRequest api_request.ChangePasswordRequest, userName string, ctx context.Context) error
	GetSalaryInformation(userName string, month string, year string, ctx context.Context) ([]*api_response.SalaryAPIResponse, error)
	ModifySalaryConfiguration(userSalaryInfo api_request.ModifySalaryConfRequest, ctx context.Context) error
	ImportStudentsByExcel(file multipart.File, courseId string, ctx context.Context) error
	ModifyUserService(rq api_request.ModifyUserInformationRequest, userId string, ctx context.Context) error
	InsertOneStudentService(request api_request.NewStudentRequest, courseId string, ctx context.Context) error
	GetCourseExistenceById(courseId string, ctx context.Context) error
	GetAllUserByJobPosition(jobPos string, ctx context.Context) ([]*UserResponse, error)
	GetStudentByCourseId(courseId string, ctx context.Context) ([]api_response.StudentResponse, error)
	AddStudentAttendanceService(rq api_request.StudentAttendanceRequest, ctx context.Context) error
	GetCourseSessionsService(courseId string, ctx context.Context) ([]api_response.StudentAttendanceScheduleResponse, error)
	UpdateStudentAttendanceService(rq api_request.StudentAttendanceRequest, ctx context.Context) error
	GetAllInChargeCourse(username string, ctx context.Context) ([]api_response.CourseResponse, error)
	CheckInWorkerAttendanceService(rq api_request.CheckInAttendanceWorkerRequest, userId string, ctx context.Context) error
	CheckEmailExistenceService(email string, ctx context.Context) (bool, error)
	PostNewForgotPasswordCode(email string, ctx context.Context) (*int, error)
	CheckFitDigitCode(email string, code int, ctx context.Context) (*bool, error)
	UpdateNewPasswordInfo(newPassword string, email string, ctx context.Context) (*api_response.UserDto, error)
	DeleteStudentService(rq api_request.DeleteStudentRequest, ctx context.Context) error
	ModifyStudentInformation(rq api_request.ModifyStudentRequest, ctx context.Context) error
	InsertNewUserByJobPosition(rq api_request.NewUserAddedByAdmin, ctx context.Context) error
	GetStudentPaymentStatusByCourseIdService(courseId string, ctx context.Context) ([]response.PaymentStatusByCourseIdResponse, error)
}

type UserStore interface {
	GetByUserNameStore(userName string, email string, userId string, ctx context.Context) (UserResponse, error)
	InsertNewUserStore(newUser user.UserEntity, ctx context.Context) error
	UpdateNewPassword(newPassword []byte, userName string) error
	GetSalaryReportStore(userName string, month string, year string, ctx context.Context) ([]salary.SalaryEntity, error)
	ModifySalaryConfigurationStore(userId string, userSalaryInfo []api_request.SalaryConfiguration, ctx context.Context) error
	InsertStudentStore(data []student.EntityStudent, courseId string, ctx context.Context) error
	ModifyUserInformationStore(rq api_request.ModifyUserInformationRequest, userId string, ctx context.Context) error
	InsertOneStudentStore(rq api_request.NewStudentRequest, courseId string, ctx context.Context) error
	CheckCourseExistence(courseId string, ctx context.Context) error
	GetUserByJobPosition(jobPos string, ctx context.Context) ([]user.UserEntity, error)
	GetUserSalaryConfigStore(userId string, ctx context.Context) ([]api_response.SalaryConfig, error)
	GetStudentByCourseIdStore(courseId string, ctx context.Context) ([]student.EntityStudent, error)
	AddStudentAttendanceStore(rq student.StudentAttendanceEntity, ctx context.Context) error
	GetScheduleByCourseIdStore(studentList []student.EntityStudent, courseId string, ctx context.Context) ([]api_response.StudentAttendanceScheduleResponse, error)
	UpdateStudentAttendanceStore(rq api_request.StudentAttendanceRequest, ctx context.Context) error
	GetUserCourseInChargeStore(username string, ctx context.Context) ([]course_class.CourseEntity, error)
	CheckInWorkerAttendanceStore(rq api_request.CheckInAttendanceWorkerRequest, userId string, ctx context.Context) error
	CheckEmailExistenceStore(email string, ctx context.Context) (bool, error)
	PostNewForgotPasswordCodeStore(email string, digitCode int, ctx context.Context) error
	CheckFitDigitCodeStore(email string, code int, ctx context.Context) (*bool, error)
	UpdateNewPasswordInfoStore(newPassword string, email string, ctx context.Context) (*user.UserEntity, error)
	DeleteStudentInCourseStore(rq api_request.DeleteStudentRequest, ctx context.Context) error
	ModifyStudentInformationStore(rq api_request.ModifyStudentRequest, ctx context.Context) error
	GetCourseManagerEntityByCourseId(courseId string, ctx context.Context) ([]course_class.CourseManagerEntity, *int, error)
}
