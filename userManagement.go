package education_website

import "context"

type UserRequest struct {
}

type UserResponse struct {
	UserId      string `json:"userId"`
	UserName    string `json:"userName"`
	DOB         string `json:"DOB"`
	JobPosition string `json:"jobPosition"`
	StartDate   string `json:"startDate"`
}

type UserService interface {
	GetByUserName(userName string, ctx context.Context) (*UserResponse, error)
}

type UserStore interface {
	GetByUserNameStore(userName string, ctx context.Context) (UserResponse, error)
}
