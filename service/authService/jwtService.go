package authService

import (
	batman "education-website"
	"education-website/entity/user"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var tokenDuration = 24 * time.Hour

type jwtService struct {
	secretKey string
}

type JwtServiceCfg struct {
	SecretKey string
}

func NewJwtService(jwtServiceCfg JwtServiceCfg) batman.JwtService {
	return &jwtService{
		secretKey: jwtServiceCfg.SecretKey,
	}
}

func (j jwtService) GenerateToken(userEntity user.UserEntity) string {
	claims := &jwt.MapClaims{
		"sub":  userEntity.UserName,
		"role": userEntity.Role,
		"exp":  time.Now().Add(tokenDuration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}
