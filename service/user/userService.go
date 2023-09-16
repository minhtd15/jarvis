package user

import (
	"context"
	batman "education-website"
	api_request "education-website/api/request"
	log "github.com/sirupsen/logrus"
)

type userService struct {
	userStore batman.UserStore
}

type UserServiceCfg struct {
	UserStore batman.UserStore
}

func NewUserService(userServiceCfg UserServiceCfg) batman.UserService {
	return &userService{
		userStore: userServiceCfg.UserStore,
	}
}

func (u userService) GetByUserName(userName string, ctx context.Context) (*batman.UserResponse, error) {
	log.Infof("Get user information by UserName")

	// init store user in here
	rs, err := u.userStore.GetByUserNameStore(userName, ctx)
	if err != nil {
		log.WithError(err).Errorf("Error getting user information from database")
		return nil, err
	}

	return &rs, nil
}

func (u userService) GetUserNamePassword(userLoginInfo api_request.LoginRequest, ctx context.Context) (*batman.UserResponse, error) {
	log.Infof("verify user information after user press sign in button")

	// init store user here
	rs, err := u.userStore.GetByUserNameStore(userLoginInfo.UserName, ctx)
	if err != nil {
		log.WithError(err).Errorf("Error getting user information from database")
		return nil, err
	}

	return &rs, nil
}
