package education_website

import (
	"context"
	api_request "education-website/api/request"
	api_response "education-website/api/response"
	"education-website/entity/salary"
	"education-website/entity/user"
)

type UserRequest struct {
	Id       int64  `json:"id"` // id = service_name + uuid + date
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
}

type UserService interface {
	GetByUserName(userName string, email string, ctx context.Context) (*UserResponse, error)
	GetUserNamePassword(userLoginInfo api_request.LoginRequest, ctx context.Context) (*UserResponse, error)
	InsertNewUser(userRegisterInfo api_request.RegisterRequest, ctx context.Context) error
	ChangePassword(changePasswordRequest api_request.ChangePasswordRequest, userName string, ctx context.Context) error
	GetSalaryInformation(userName string, month string, year string) (*api_response.SalaryResponse, error)
}

type UserStore interface {
	GetByUserNameStore(userName string, email string, ctx context.Context) (UserResponse, error)
	InsertNewUserStore(userRegisterInfo user.UserEntity, ctx context.Context) error
	UpdateNewPassword(newPassword []byte, userName string) error
	GetSalaryReportStore(userName string, month string, year string) (salary.SalaryEntity, error)
}
