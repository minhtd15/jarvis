package user

import (
	"context"
	batman "education-website"
	api_request "education-website/api/request"
	api_response "education-website/api/response"
	"education-website/entity/student"
	"education-website/entity/user"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
	"mime/multipart"
	"net/smtp"
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
				UserId:      repoRes.UserId,
				UserName:    repoRes.UserName,
				FullName:    repoRes.FullName,
				Gender:      repoRes.Gender,
				JobPosition: repoRes.JobPosition,
			}

			userSalaryConfig, err := u.userStore.GetUserSalaryConfigStore(x.UserId, ctx)
			if err != nil {
				log.WithError(err).Errorf("Error getting user salary config")
				return nil, err
			}

			salaryInfo := make([]api_response.SalaryInformation, 0)
			info := api_response.SalaryInformation{
				PayrollId: repoRes.PayrollId,
				WorkDays:  repoRes.TotalWorkDates,
				Amount:    repoRes.TotalSalary,
			}
			salaryInfo = append(salaryInfo, info)
			x.Salary = salaryInfo
			x.SalaryConfig = userSalaryConfig
			res = append(res, &x)
		} else {
			for i := range res {
				if res[i].UserName == repoRes.UserName {
					salary := res[i].Salary
					info := api_response.SalaryInformation{
						PayrollId: repoRes.PayrollId,
						WorkDays:  repoRes.TotalWorkDates,
						Amount:    repoRes.TotalSalary,
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
			file.SetCellValue("Sheet1", col+strconv.Itoa(row+1), salaryInfo.WorkDays)
			file.SetCellValue("Sheet1", col+strconv.Itoa(row+3), salaryInfo.Amount)
		}
	}

	return file, nil
}

func (u userService) ImportStudentsByExcel(file multipart.File, courseId string, ctx context.Context) error {
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

		dobTmp, err := time.Parse("2006-01-02", dob)
		if err != nil {
			log.WithError(err).Errorf("Unable to parse time")
			return err
		}
		rowData := student.EntityStudent{
			Name:    studentName,
			DOB:     dobTmp,
			Email:   email,
			PhoneNo: phone,
		}

		studentData = append(studentData, rowData)
	}

	err = u.userStore.InsertStudentStore(studentData, courseId, ctx)
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

func (u userService) InsertOneStudentService(request api_request.NewStudentRequest, courseId string, ctx context.Context) error {
	return u.userStore.InsertOneStudentStore(request, courseId, ctx)
}

func SendDailyEmail() {
	// Dữ liệu cần hiển thị trong email (ví dụ)
	data := "Dữ liệu của bạn: <strong>Thông tin lịch làm ngày hôm nay</strong>"

	// Định dạng nội dung email với dữ liệu
	emailBody := fmt.Sprintf("<html><body>%s</body></html>", data)

	// Gửi email
	err := sendEmail("cthanhnguyen03@gmail.com", "Subject: Daily Schedule", emailBody)
	if err != nil {
		fmt.Println("Error sending email:", err)
	}
}

func sendEmail(recipient, subject, body string) error {
	// Địa chỉ email và mật khẩu của người gửi
	from := "ducminhtong1510@gmail.com"
	password := "hiks irqs gwyz eygn"

	// Địa chỉ SMTP server và cổng
	smtpServer := "smtp.gmail.com"
	smtpPort := 587

	// Tạo một cấu trúc đại diện cho thông tin đăng nhập
	auth := smtp.PlainAuth("", from, password, smtpServer)

	// Định dạng nội dung email dưới dạng HTML
	message := fmt.Sprintf("Subject: %s\r\n", subject)
	message += "MIME-version: 1.0;\r\n"
	message += "Content-Type: text/html; charset=\"UTF-8\";\r\n"
	message += "\r\n" + body

	// Gửi email
	err := smtp.SendMail(fmt.Sprintf("%s:%d", smtpServer, smtpPort), auth, from, []string{recipient}, []byte(message))
	if err != nil {
		return err
	}

	return nil
}

func (u userService) GetCourseExistenceById(courseId string, ctx context.Context) error {
	return u.userStore.CheckCourseExistence(courseId, ctx)
}

func (u userService) GetAllUserByJobPosition(jobPos string, ctx context.Context) ([]*batman.UserResponse, error) {
	log.Infof("Start service get all %s user", jobPos)
	userList, err := u.userStore.GetUserByJobPosition(jobPos, ctx)
	if err != nil {
		log.WithError(err).Errorf("Error get all %s user from db", jobPos)
		return nil, err
	}

	var rs []*batman.UserResponse
	for _, v := range userList {
		var tmp = batman.UserResponse{
			UserId:      v.UserId,
			UserName:    v.UserName,
			DOB:         v.DOB,
			Email:       v.Email,
			JobPosition: v.JobPosition,
			Role:        v.Role,
			StartDate:   v.StartingDate,
			FullName:    v.FullName,
		}
		rs = append(rs, &tmp)
	}

	log.Infof("Done service get user by job position")
	return rs, nil
}

func (u userService) GetStudentByCourseId(courseId string, ctx context.Context) ([]api_response.StudentResponse, error) {
	log.Infof("Get student info by course ID")

	studentEntity, err := u.userStore.GetStudentByCourseIdStore(courseId, ctx)
	if err != nil {
		log.WithError(err).Errorf("Error get all student in course %s from db", courseId)
		return nil, err
	}
	var rs []api_response.StudentResponse
	for _, v := range studentEntity {
		tmp := api_response.StudentResponse{
			StudentId:   v.Id,
			StudentName: v.Name,
			Email:       v.Email,
			PhoneNumber: v.PhoneNo,
			Dob:         v.DOB,
		}
		rs = append(rs, tmp)
	}
	return rs, nil
}

func (u userService) AddStudentAttendanceService(rq api_request.StudentAttendanceRequest, ctx context.Context) error {
	log.Infof("Service add student request for student %s in classID %s and status is %s", rq.StudentId, rq.ClassId, rq.Status)

	entity := student.StudentAttendanceEntity{
		StudentId: rq.StudentId,
		ClassId:   rq.ClassId,
		Status:    rq.Status,
	}

	err := u.userStore.AddStudentAttendanceStore(entity, ctx)
	if err != nil {
		log.WithError(err).Errorf("Error service add student attendance in store/database")
		return err
	}

	return nil
}

func (u userService) GetCourseSessionsService(courseId string, ctx context.Context) ([]api_response.StudentAttendanceScheduleResponse, error) {
	log.Infof("First need to get the student list from course %s", courseId)

	studentList, err := u.userStore.GetStudentByCourseIdStore(courseId, ctx)
	if err != nil {
		log.WithError(err).Errorf("Unable to get student list")
		return nil, err
	}

	rs, err := u.userStore.GetScheduleByCourseIdStore(studentList, courseId, ctx)
	if err != nil {
		log.WithError(err).Errorf("Unable to get student attendance along with schedule")
		return nil, err
	}

	return rs, nil
}

func (u userService) UpdateStudentAttendanceService(rq api_request.StudentAttendanceRequest, ctx context.Context) error {
	return u.userStore.UpdateStudentAttendanceStore(rq, ctx)
}

func (u userService) GetAllInChargeCourse(username string, ctx context.Context) ([]api_response.CourseResponse, error) {
	log.Infof("Start to get course which user %s in charge", username)

	rs, err := u.userStore.GetUserCourseInChargeStore(username, ctx)
	if err != nil {
		log.WithError(err).Errorf("Unable to get course that user in charge store")
		return nil, err
	}

	var result []api_response.CourseResponse

	for _, v := range rs {
		tmp := api_response.CourseResponse{
			CourseId:   strconv.Itoa(v.CourseId),
			CourseName: v.CourseName,
			Room:       strconv.Itoa(v.Room),
			StartDate:  v.StartDate.String,
			EndDate:    v.EndDate.String,
			StudyDays:  v.StudyDays,
			Location:   v.Location.String,
		}
		result = append(result, tmp)
	}

	return result, nil
}

func (u userService) CheckInWorkerAttendanceService(rq api_request.CheckInAttendanceWorkerRequest, userId string, ctx context.Context) error {
	log.Infof("Start to get count attendance service for worker %v", rq)
	return u.userStore.CheckInWorkerAttendanceStore(rq, userId, ctx)
}
