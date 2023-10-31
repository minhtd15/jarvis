package education_website

import (
	"context"
	api_request "education-website/api/request"
	api_response "education-website/api/response"
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
	UserId      string
	UserName    string
	DOB         string
	Email       string
	JobPosition string
	Role        string
	StartDate   string
	Password    string
	FullName    string
	Gender      string
}

type UserService interface {
	GetByUserName(userName string, email string, userId string, ctx context.Context) (*UserResponse, error)
	GetUserNamePassword(userLoginInfo api_request.LoginRequest, ctx context.Context) (*UserResponse, error)
	InsertNewUser(userRegisterInfo api_request.RegisterRequest, ctx context.Context) error
	ChangePassword(changePasswordRequest api_request.ChangePasswordRequest, userName string, ctx context.Context) error
	GetSalaryInformation(userName string, month string, year string, ctx context.Context) ([]*api_response.SalaryAPIResponse, error)
	ModifySalaryConfiguration(userSalaryInfo api_request.ModifySalaryConfRequest, ctx context.Context) error
	ImportStudentsByExcel(file multipart.File, ctx context.Context) error
}

type UserStore interface {
	GetByUserNameStore(userName string, email string, userId string, ctx context.Context) (UserResponse, error)
	InsertNewUserStore(userRegisterInfo user.UserEntity, ctx context.Context) error
	UpdateNewPassword(newPassword []byte, userName string) error
	GetSalaryReportStore(userName string, month string, year string, ctx context.Context) ([]salary.SalaryEntity, error)
	ModifySalaryConfigurationStore(userId string, userSalaryInfo []api_request.SalaryConfiguration, ctx context.Context) error
	InsertStudentStore(data []student.EntityStudent, ctx context.Context) error
}
