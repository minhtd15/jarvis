package authService

import (
	batman "education-website"
	"education-website/entity/user"
	"errors"
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

func NewPayload(username string, userId string, role string, duration time.Duration) (*batman.Payload, error) {
	payload := &batman.Payload{
		Username:  username,
		UserId:    userId,
		Role:      role,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (j jwtService) GenerateToken(userEntity user.UserEntity) string {
	payload, err := NewPayload(userEntity.UserName, userEntity.UserId, userEntity.Role, tokenDuration)
	if err != nil {
		return ""
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (j *jwtService) ValidateToken(tokenString string) (*batman.Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("token has expired")
		}
		return []byte(j.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(tokenString, &batman.Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, errors.New("token has expired")) {
			return nil, errors.New("token has expired")
		}
		return nil, errors.New("token has expired")
	}

	payload, ok := jwtToken.Claims.(*batman.Payload)
	if !ok {
		return nil, errors.New("token has expired")
	}

	return payload, nil
}
