package login

import (
	batman "education-website"
)

type userService struct {
	userStore batman.UserStore
}

type UserServiceCfg struct {
	UserStore batman.UserStore
}

func NewUserService(userServiceCfg UserServiceCfg) *userService {
	return &userService{
		userStore: userServiceCfg.UserStore,
	}
}
