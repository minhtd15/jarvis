package service

import (
	"context"
	batman "education-website"
	log "github.com/sirupsen/logrus"
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

// hello
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
