package authService

import (
	batman "education-website"
	"education-website/api/request"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
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
	if err := bcrypt.CompareHashAndPassword([]byte(userEntity.Password), []byte(userLoginRequest.Password)); err != nil {
		// Trả về một lỗi hoặc mã lỗi xác định
		return nil, errors.New("Invalid username or password")
	}

	if userLoginRequest.Email == userEntity.Email {
		userInfo := map[string]interface{}{
			"username": userLoginRequest.Email,
			"role":     "user",
		}
		return userInfo, nil
	}

	// invalid username or password, return error
	return nil, fmt.Errorf("Invalid username or password")
}
