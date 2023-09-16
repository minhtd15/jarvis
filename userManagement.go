package education_website

import (
	"context"
	api_request "education-website/api/request"
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
	GetByUserName(userName string, ctx context.Context) (*UserResponse, error)
	GetUserNamePassword(userLoginInfo api_request.LoginRequest, ctx context.Context) (*UserResponse, error)
}

type UserStore interface {
	GetByUserNameStore(userName string, ctx context.Context) (UserResponse, error)
}
