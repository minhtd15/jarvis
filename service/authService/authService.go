package authService

import (
	batman "education-website"
	"education-website/api/request"
	"fmt"
)

type authService struct {
	jwtService batman.JwtService
}

type AuthServiceCfg struct {
	JwtService batman.JwtService
}

func NewAuthService(authServiceCfg AuthServiceCfg) *authService {
	return &authService{
		jwtService: authServiceCfg.JwtService,
	}
}

func (a authService) VerifyUser(userLoginRequest request.LoginRequest, userEntity batman.UserResponse) (interface{}, error) {
	// this service is used to check the login password similar with the password on database
	if userLoginRequest.UserName == userEntity.UserName && userLoginRequest.Password == userEntity.Password {
		// Trả về thông tin người dùng sau khi xác thực thành công.
		userInfo := map[string]interface{}{
			"username": userLoginRequest.UserName,
			"role":     "user",
		}
		return userInfo, nil
	}

	// invalid username or password, return error
	return nil, fmt.Errorf("Invalid username or password")
}
