package user

import (
	"context"
	batman "education-website"
	api_request "education-website/api/request"
	api_response "education-website/api/response"
	"education-website/entity/student"
	"education-website/entity/user"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
	"mime/multipart"
	"strconv"
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

	if err := bcrypt.CompareHashAndPassword([]byte(oldPasswordEntity.Password), []byte(oldPassword)); err != nil {
		// Trả về một lỗi hoặc mã lỗi xác định
		log.WithError(err).Errorf("Invalid old password, %s", err)
		return errors.New("Invalid username or password")
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
				UserName:    repoRes.UserName,
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

func ExportToExcel(data []api_response.SalaryAPIResponse) (*excelize.File, error) {
	file := excelize.NewFile()

	for index, item := range data {
		// Điều chỉnh index để bắt đầu từ dòng 2 (dòng tiêu đề ở dòng 1)
		row := index + 2

		// Ghi dữ liệu vào các ô trong tệp Excel
		file.SetCellValue("Sheet1", "A"+strconv.Itoa(row), item.UserName)
		file.SetCellValue("Sheet1", "B"+strconv.Itoa(row), item.FullName)
		file.SetCellValue("Sheet1", "C"+strconv.Itoa(row), item.Gender)
		file.SetCellValue("Sheet1", "D"+strconv.Itoa(row), item.JobPosition)

		// Ghi dữ liệu SalaryInformation vào các cột tương ứng
		for i, salaryInfo := range item.Salary {
			col := string('E' + i)
			file.SetCellValue("Sheet1", col+strconv.Itoa(row), salaryInfo.CourseType)
			file.SetCellValue("Sheet1", col+strconv.Itoa(row+1), salaryInfo.WorkDays)
			file.SetCellValue("Sheet1", col+strconv.Itoa(row+2), salaryInfo.PriceEach)
			file.SetCellValue("Sheet1", col+strconv.Itoa(row+3), salaryInfo.Amount)
		}
	}

	return file, nil
}

func (u userService) ImportStudentsByExcel(file multipart.File, ctx context.Context) error {
	// Read the Excel file
	f, err := excelize.OpenReader(file)
	if err != nil {
		log.WithError(err).Errorf("Unable to open excel file")
		return err
	}

	// Get values from the specified sheet and columns
	sheetName := "Student"
	columnName := "A"    // Column for studentName
	dobColumnName := "B" // Column for DOB
	emailColumn := "C"
	phoneColumn := "D"

	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.WithError(err).Errorf("Error getting row from excel file")
		return err
	}

	var studentData []student.EntityStudent

	// Iterate through rows and extract data
	for i, row := range rows {
		// skip the first row as the headers of table
		if i == 0 {
			continue
		}

		studentName := u.rowToColumnValue(row, columnName)
		dob := u.rowToColumnValue(row, dobColumnName)
		email := u.rowToColumnValue(row, emailColumn)
		phone := u.rowToColumnValue(row, phoneColumn)

		rowData := student.EntityStudent{
			Name:    studentName,
			DOB:     dob,
			Email:   email,
			PhoneNo: phone,
		}

		studentData = append(studentData, rowData)
	}

	err = u.userStore.InsertStudentStore(studentData, ctx)
	if err != nil {
		log.WithError(err).Errorf("Error inserting student data to database")
		return err
	}

	return nil
}

func (u userService) rowToColumnValue(row []string, column string) string {
	columnIndex, err := excelize.ColumnNameToNumber(column)
	if err != nil {
		log.Fatal(err)
	}

	if columnIndex <= len(row) {
		return row[columnIndex-1]
	}

	return "" // Return empty string if column index is out of range
}

func (u userService) ModifyUserService(rq api_request.ModifyUserInformationRequest, userId string, ctx context.Context) error {
	return u.userStore.ModifyUserInformationStore(rq, userId, ctx)
}
