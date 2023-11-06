package education_website

import (
	api_request "education-website/api/request"
	"education-website/entity/user"
	"errors"
	"time"
)

type Payload struct {
	Username     string    `json:"username"`
	Role         string    `json:"role"`
	UserId       string    `json:"user_id"`
	UserFullName string    `json:"user_fullname"`
	IssuedAt     time.Time `json:"issued_at"`
	ExpiredAt    time.Time `json:"expired_at"`
}

func (p Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return errors.New("Token has expired")
	}
	return nil
}

type AuthService interface {
	VerifyUser(userLoginRequest api_request.LoginRequest, userEntity UserResponse) (interface{}, error)
}

type JwtService interface {
	GenerateToken(userEntity user.UserEntity) string
	ValidateToken(token string) (*Payload, error)
}
