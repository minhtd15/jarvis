package user

import (
	"context"
	"database/sql"
	batman "education-website"
	api_request "education-website/api/request"
	"education-website/client"
	"education-website/entity/user"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type userService struct {
	userStore   batman.UserStore
	flashClient client.FlashClient
}

type UserServiceCfg struct {
	UserStore   batman.UserStore
	FlashClient client.FlashClient
}

func NewUserService(userServiceCfg UserServiceCfg) batman.UserService {
	return &userService{
		userStore:   userServiceCfg.UserStore,
		flashClient: userServiceCfg.FlashClient,
	}
}

func (u userService) GetByUserName(userName string, email string, userId string, ctx context.Context) (*batman.UserResponse, error) {
	log.Infof("Get user information by UserName")

	// init store user in here
	rs, err := u.userStore.GetByUserNameStore(userName, email, userId, ctx)
	if errors.Is(err, sql.ErrNoRows) {
		log.WithError(err).Errorf("No user found")
		return nil, err
	}
	if err != nil {
		log.WithError(err).Errorf("Error getting user information from database")
		return nil, err
	}

	return &rs, nil
}

func (u userService) GetUserNamePassword(userLoginInfo api_request.LoginRequest, ctx context.Context) (*batman.UserResponse, error) {
	log.Infof("verify user information after user press sign in button")

	// init store user here
	rs, err := u.userStore.GetByUserNameStore("", userLoginInfo.Email, "", ctx)
	if err != nil {
		log.WithError(err).Errorf("Error getting user information from database")
		return nil, err
	}

	return &rs, nil
}

func (u userService) InsertNewUser(userRegisterInfo api_request.RegisterRequest, ctx context.Context) (string, error) {
	log.Infof("insert new user to database when user register for a new account")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRegisterInfo.Password), 12)
	if err != nil {
		log.WithError(err).Errorf("Error encrypt password")
		return "0", err
	}
	newUser := user.UserEntity{
		UserId:       u.GenerateUserId(),
		UserName:     userRegisterInfo.UserName,
		Email:        userRegisterInfo.Email,
		Role:         "user",
		DOB:          userRegisterInfo.DOB,
		StartingDate: time.Now().Format("2006-01-02"),
		JobPosition:  "Undefined",
		Password:     string(hashedPassword),
		Gender:       userRegisterInfo.Gender,
		FullName:     userRegisterInfo.FullName,
	}

	err = u.userStore.InsertNewUserStore(newUser, ctx)
	if err != nil {
		log.WithError(err).Errorf("Error insert new user to DB")
		return "0", nil
	}

	return newUser.UserId, nil
}

func (u userService) GenerateUserId() string {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	currentTime = strings.ReplaceAll(currentTime, "-", "")
	currentTime = strings.ReplaceAll(currentTime, ":", "")
	currentTime = strings.ReplaceAll(currentTime, " ", "")
	return currentTime
}

func (u userService) ChangePassword(changePasswordRequest api_request.ChangePasswordRequest, userName string, ctx context.Context) error {
	log.Infof("Start to change password")

	err := u.VerifyChangePassword(changePasswordRequest.OldPassword, userName, ctx)
	if err != nil {
		log.WithError(err).Errorf("Unable to verify changing password")
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(changePasswordRequest.NewPassword), 12)
	err = u.userStore.UpdateNewPassword(hashedPassword, userName)
	if err != nil {
		log.WithError(err).Errorf("Unable tp update new password")
		return err
	}
	return nil
}

func (u userService) VerifyChangePassword(oldPassword string, userName string, ctx context.Context) error {
	log.Info("Compare password in database")

	oldPasswordEntity, err := u.userStore.GetByUserNameStore(userName, "", "", ctx)
	if err != nil {
		log.WithError(err).Errorf("unable to get password from database")
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(oldPasswordEntity.Password), []byte(oldPassword)); err != nil {
		// Trả về một lỗi hoặc mã lỗi xác định
		log.WithError(err).Errorf("Invalid old password, %s", err)
		return errors.New("Invalid username or password")
	}

	return nil
}
