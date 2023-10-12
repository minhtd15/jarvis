package user

import (
	"context"
	batman "education-website"
	api_request "education-website/api/request"
	api_response "education-website/api/response"
	"education-website/entity/user"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
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

func (u userService) GetByUserName(userName string, email string, userId string, ctx context.Context) (*batman.UserResponse, error) {
	log.Infof("Get user information by UserName")

	// init store user in here
	rs, err := u.userStore.GetByUserNameStore(userName, email, userId, ctx)
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

func (u userService) InsertNewUser(userRegisterInfo api_request.RegisterRequest, ctx context.Context) error {
	log.Infof("insert new user to database when user register for a new account")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRegisterInfo.Password), 12)
	if err != nil {
		log.WithError(err).Errorf("Error encrypt password")
		return err
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
		return err
	}

	return nil
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

	// compare old and new password
	err = bcrypt.CompareHashAndPassword([]byte(oldPassword), []byte(oldPasswordEntity.Password))
	if err != nil {
		log.WithError(err).Errorf("Invalid old password")
		return fmt.Errorf("invalid old password")
	}

	return nil
}

func (u userService) GetSalaryInformation(userName string, month string, year string, ctx context.Context) ([]*api_response.SalaryAPIResponse, error) {
	log.Infof("Get salary information for user %s", userName)

	userSalaryReport, err := u.userStore.GetSalaryReportStore(userName, month, year, ctx)
	if err != nil {
		log.WithError(err).Errorf("Unable to get salary report for user %s", userName)
		return nil, err
	}

	res := make([]*api_response.SalaryAPIResponse, 0) // This is response body
	check := make(map[string]bool)                    // variable is used to check whether userId exists in response
	for _, repoRes := range userSalaryReport {
		if !check[repoRes.UserName] {
			check[repoRes.UserName] = true
			x := api_response.SalaryAPIResponse{
				UserName:    userName,
				FullName:    repoRes.FullName,
				Gender:      repoRes.Gender,
				JobPosition: repoRes.JobPosition,
			}

			salaryInfo := make([]api_response.SalaryInformation, 0)
			info := api_response.SalaryInformation{
				CourseType: repoRes.TypeWork,
				WorkDays:   repoRes.TotalWorkDates,
				PriceEach:  repoRes.PayrollPerSessions,
				Amount:     repoRes.TotalSalary,
			}
			salaryInfo = append(salaryInfo, info)
			x.Salary = salaryInfo
			res = append(res, &x)
		} else {
			for i := range res {
				if res[i].UserName == repoRes.UserName {
					salary := res[i].Salary
					info := api_response.SalaryInformation{
						CourseType: repoRes.TypeWork,
						WorkDays:   repoRes.TotalWorkDates,
						PriceEach:  repoRes.PayrollPerSessions,
						Amount:     repoRes.TotalSalary,
					}
					salary = append(salary, info)
					res[i].Salary = salary
				}
			}
		}
	}

	return res, nil

}

func (u userService) ModifySalaryConfiguration(userSalaryInfo api_request.ModifySalaryConfRequest, ctx context.Context) error {
	return u.userStore.ModifySalaryConfigurationStore(userSalaryInfo.UserId, userSalaryInfo.NewSalaryList, ctx)
}
