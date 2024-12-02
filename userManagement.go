package education_website

import (
	"context"
	api_request "education-website/api/request"
	"education-website/entity/salary"
	"education-website/entity/user"
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
}

type UserStore interface {
	GetByUserNameStore(userName string, email string, userId string, ctx context.Context) (UserResponse, error)
	InsertNewUserStore(newUser user.UserEntity, ctx context.Context) error
	UpdateNewPassword(newPassword []byte, userName string) error
	GetSalaryReportStore(userName string, month string, year string, ctx context.Context) ([]salary.SalaryEntity, error)
}
