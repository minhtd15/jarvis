package education_website

import (
	api_request "education-website/api/request"
	"education-website/entity/user"
)

type jwtResponse struct {
}

type AuthService interface {
	VerifyUser(userLoginRequest api_request.LoginRequest, userEntity UserResponse) (interface{}, error)
}

type JwtService interface {
	GenerateToken(userEntity user.UserEntity) string
}
