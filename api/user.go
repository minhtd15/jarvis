package api

import (
	"context"
	batman "education-website"
	api_request "education-website/api/request"
	api_response "education-website/api/response"
	"education-website/commonconstant"
	"education-website/service/user"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	_ "github.com/xuri/excelize/v2"
	"go.elastic.co/apm"
	"io/ioutil"
	"net/http"
)

func handlerUserAccount(w http.ResponseWriter, r *http.Request) {

	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD", "handleRetryUserAccount")
	logger.Infof("Handle user account")

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Warningf("Error when reading from request")
		http.Error(w, "Invalid format", 252001)
		return
	}

	json.NewDecoder(r.Body)
	defer r.Body.Close()
	var userRequest batman.UserRequest
	err = json.Unmarshal(bodyBytes, &userRequest)
	if err != nil {
		log.WithError(err).Warningf("Error when unmarshaling data from request")
		http.Error(w, "Status bad Request", http.StatusBadRequest) // Return a 400 Bad Request error
		return
	}

	// Make sure userService is initialized and not nil
	if userService == nil {
		logger.Errorf("userService is nil")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// get user's information from database
	log.Infof("Start to get user's information from database")
	user, err := userService.GetByUserName(userRequest.UserName, userRequest.Email, userRequest.Id, ctx)
	if err != nil {
		log.WithError(err).Errorf("Username %s does not exist", user.UserName)
		return
	}

	log.Infof("Successful get the user data")
	respondWithJSON(w, http.StatusOK, batman.CommonResponse{Status: "SUCCESS", Descrition: "Success getting the user data"})
}

func handlerSalaryInformation(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD", "handle get salary information")
	logger.Infof("Handle user account")

	keys := r.URL.Query()
	userSearchName := keys.Get("username") // this variable is used for searching engine for leader
	month := keys.Get("month")
	year := keys.Get("year")
	role, ok := r.Context().Value("role").(string)
	userName, ok := r.Context().Value("user_fullname").(string)
	if !ok {
		http.Error(w, "Unable to get role/userName from token", http.StatusUnauthorized)
		return
	}

	salaryRequest := api_request.SalaryRequest{
		Month: month,
		Year:  year,
	}

	if role == "user" {
		userSalaryReport, err := userService.GetSalaryInformation(userName, salaryRequest.Month, salaryRequest.Year, ctx)
		if err != nil {
			log.WithError(err).Warningf("Error getting salary information from Salary view for user %s", userName)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"message": "Successful getting user salary information",
			"data":    userSalaryReport,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	} else {
		if userSearchName == "" {
			userSearchName = "undefined"
			userSalaryReport, err := userService.GetSalaryInformation(userSearchName, salaryRequest.Month, salaryRequest.Year, ctx)
			if err != nil {
				log.WithError(err).Warningf("Error getting salary information from Salary view for leader %s", userName)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			response := map[string]interface{}{
				"message": "Successful getting all user salary information for leader",
				"data":    userSalaryReport,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		} else {
			userSalaryReport, err := userService.GetSalaryInformation(userSearchName, salaryRequest.Month, salaryRequest.Year, ctx)
			if err != nil {
				log.WithError(err).Warningf("Error getting salary information from Salary view for leader %s", userName)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			response := map[string]interface{}{
				"message": "Successful getting user information: " + userSearchName,
				"data":    userSalaryReport,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}

	}
}

func handleModifySalaryConfiguration(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD", "handle get salary information")
	logger.Infof("Handle user account")

	//keys := r.URL.Query() // this variable is used for searching engine for leader
	role, ok := r.Context().Value("role").(string)
	if !ok {
		http.Error(w, "Unable to get role/userName from token", http.StatusUnauthorized)
		return
	}

	if role == "user" {
		response := map[string]interface{}{
			"message": "You are not allowed to access to this function",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Warningf("Error when reading from request")
		http.Error(w, "Invalid format", 252001)
		return
	}

	json.NewDecoder(r.Body)
	defer r.Body.Close()
	var newSalaryInfo api_request.ModifySalaryConfRequest
	err = json.Unmarshal(bodyBytes, &newSalaryInfo)
	if err != nil {
		log.WithError(err).Warningf("Error when unmarshaling data from request")
		http.Error(w, "Status internal server error", http.StatusInternalServerError) // Return a 400 Bad Request error
		return
	}

	err = userService.ModifySalaryConfiguration(newSalaryInfo, ctx)
	if err != nil {
		log.WithError(err).Errorf("Error modify salary configuration for user")
		http.Error(w, "Status internal server error", http.StatusInternalServerError)
		return
	}
}

func handleExcelSalary(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD", "handle get salary information")
	logger.Infof("Handle user account")

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Warningf("Error when reading from request")
		http.Error(w, "Invalid format", 252001)
		return
	}

	json.NewDecoder(r.Body)
	defer r.Body.Close()
	var data []api_response.SalaryAPIResponse
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		log.WithError(err).Warningf("Error when unmarshaling data from request")
		http.Error(w, "Status bad Request", http.StatusBadRequest) // Return a 400 Bad Request error
		return
	}

	excelFile, err := user.ExportToExcel(data)
	if err != nil {
		http.Error(w, "Failed to create Excel file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=output.xlsx")

	// Save the Excel file to the response writer
	if err := excelFile.Write(w); err != nil {
		fmt.Println("Error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func handleInsertStudents(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD", "handle insert students information")
	logger.Infof("Insert all students information from an excel file")

	keys := r.URL.Query()
	courseId := keys.Get("course_id")
	if courseId == "" {
		// courseId is missing, return an error
		log.Error("course_id parameter is missing")
		http.Error(w, "course_id parameter is required", http.StatusBadRequest)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB limit for the uploaded file
	if err != nil {
		log.WithError(err).Errorf("ERROR IN UPLOAD FILE: %s", err)
		http.Error(w, "Unable to parse uploaded file", http.StatusBadRequest)
		return
	}

	// Get the uploaded file
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Missing or invalid 'file' field in the request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// First need to check whether the course exists in DB or not
	err = userService.GetCourseExistenceById(courseId, ctx)
	if err != nil {
		if errors.Is(err, commonconstant.ErrCourseNotExist) {
			log.WithError(err).Errorf("Cannot find the courseId %s", courseId)
			http.Error(w, "Cannot find course according to requirement", http.StatusBadRequest)
			return
		}
		log.WithError(err).Errorf("Error checking course existence: %s", err)
		http.Error(w, "error checking course existence", http.StatusInternalServerError)
		return
	}

	err = userService.ImportStudentsByExcel(file, courseId, ctx)
	if err != nil {
		log.WithError(err).Errorf("Error insert student to db by import excel")
		http.Error(w, "Cannot insert student to db by import excel", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Successful insert students list to database by importing excel data",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleModifyUserInformation(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD", "modify user information")
	logger.Infof("Only leader and admin can modify this information")

	userId, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unable to get role/userName from token", http.StatusUnauthorized)
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Warningf("Error when reading from request")
		http.Error(w, "Invalid format", 252001)
		return
	}

	json.NewDecoder(r.Body)
	defer r.Body.Close()

	// this is the information that the user type in front end
	var rq api_request.ModifyUserInformationRequest
	err = json.Unmarshal(bodyBytes, &rq)
	if err != nil {
		log.WithError(err).Errorf("Error marshaling body to modify user information request")
		http.Error(w, "Status internal Request", http.StatusInternalServerError)
		return
	}

	err = userService.ModifyUserService(rq, userId, ctx)
	if err != nil {
		log.WithError(err).Errorf("Error modify user information request service")
		http.Error(w, "Error internal Request", http.StatusInternalServerError)
		return
	}

	log.Infof("Successful update user %s information", userId)
	response := map[string]interface{}{
		"message": "Successful insert students list to database by importing excel data",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleInsertOneNewStudent(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD", "insert one student into class")
	logger.Infof("Only leader and admin can insert student")

	role, ok := r.Context().Value("role").(string)
	if !ok {
		http.Error(w, "Unable to get role/userName from token", http.StatusUnauthorized)
		return
	}

	if role == "user" {
		response := map[string]interface{}{
			"message": "You are not allowed to this function",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(response)
	}

	//keys := r.URL.Query()
	//courseId := keys.Get("course_id")
	//if courseId == "" {
	//	// courseId is missing, return an error
	//	log.Error("course_id parameter is missing")
	//	http.Error(w, "course_id parameter is required", http.StatusBadRequest)
	//	return
	//}
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Warningf("Error when reading from request")
		http.Error(w, "Invalid format", 252001)
		return
	}

	json.NewDecoder(r.Body)
	defer r.Body.Close()

	// this is the information that the user type in front end
	var rq api_request.NewStudentRequest
	err = json.Unmarshal(bodyBytes, &rq)
	if err != nil {
		log.WithError(err).Errorf("Error marshaling body to modify user information request")
		http.Error(w, "Status internal Request", http.StatusInternalServerError)
		return
	}

	// First need to check whether the course exists in DB or not
	err = userService.GetCourseExistenceById(rq.CourseId, ctx)
	if err != nil {
		if errors.Is(err, commonconstant.ErrCourseNotExist) {
			log.WithError(err).Errorf("Cannot find the courseId %s", rq.CourseId)
			http.Error(w, "Cannot find course according to requirement", http.StatusBadRequest)
			return
		}
		log.WithError(err).Errorf("Error checking course existence: %s", err)
		http.Error(w, "error checking course existence", http.StatusInternalServerError)
		return
	}

	err = userService.InsertOneStudentService(rq, rq.CourseId, ctx)
	if err != nil {
		log.WithError(err).Errorf("Error insert student %s", rq.Name)
		http.Error(w, "Unable to insert one new student", http.StatusInternalServerError)
		return
	}

	log.Infof("Successful update student %s information", rq.Name)
	response := map[string]interface{}{
		"message": "Successful insert one student list to database",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleEmailJobs(w http.ResponseWriter, r *http.Request) {
	log.Infof("Start to send email")

	// Replace these values with your AWS credentials and SES region
	accessKey := "AKIARRJQQHPNOGS3BIUC"
	secretKey := "ROf1eSmtSasK2KBcAOrxrVSr7ocmkImI/X0aN/x+"
	region := "ap-southeast-2"

	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		log.Errorf("Error creating AWS session: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create an SES client
	svc := ses.New(sess)

	// Specify the sender email address
	sender := "ducminhtong1510@gmail.com"

	// Specify the recipient email address
	recipient := "minhtom1510@gmail.com"

	// Specify the email subject and body
	subject := "Test Email"
	body := "This is a test email sent using Amazon SES."

	// Send the email
	_, err = svc.SendEmailWithContext(context.TODO(), &ses.SendEmailInput{
		Source: aws.String(sender),
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(recipient)},
		},
		Message: &ses.Message{
			Subject: &ses.Content{
				Data: aws.String(subject),
			},
			Body: &ses.Body{
				Text: &ses.Content{
					Data: aws.String(body),
				},
			},
		},
	})
	if err != nil {
		log.Errorf("Error sending email: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Infof("Email sent successfully")
	fmt.Fprint(w, "Email sent successfully")
}

func handleGetUserByRole(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD", "get all teacher or role according to request")
	logger.Infof("get user who are teacher or TA according to request")

	keys := r.URL.Query()
	jobPos := keys.Get("job_position")
	if jobPos == "" || (jobPos != "Teacher" && jobPos != "TA") {
		// courseId is missing, return an error
		log.Error("Job position parameter is missing or job Position is not Teacher or TA")
		http.Error(w, "job position parameter is required or job position must be Teacher or TA in the right format", http.StatusBadRequest)
		return
	}

	userList, err := userService.GetAllUserByJobPosition(jobPos, ctx)
	if err != nil {
		log.WithError(err).Errorf("Error getting users who are %s", jobPos)
		http.Error(w, "Cannot get teacher/TA", http.StatusInternalServerError)
	}

	response := map[string]interface{}{
		"message": "Successful get  " + jobPos,
		"data":    userList,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleDeleteCourse(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD", "delete course according to request")
	logger.Infof("API Delete course")

	role, ok := r.Context().Value("role").(string)
	if !ok {
		http.Error(w, "Unable to get role/userName from token", http.StatusUnauthorized)
		return
	}

	if role == "user" {
		response := map[string]interface{}{
			"message": "You are not allowed to this function",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(response)
	}

	keys := r.URL.Query()
	courseId := keys.Get("course_id")
	if courseId == "" {
		// courseId is missing, return an error
		log.Error("CourseId parameter is missing")
		http.Error(w, "CourseId parameter is required ", http.StatusBadRequest)
		return
	}

	err := classService.DeleteCourseByCourseId(courseId, ctx)
	if err != nil {
		log.WithError(err).Errorf("Error delete COURSE %s", courseId)
		http.Error(w, "Error delete course. Please try again", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

//func handleDeleteUser(w http.ResponseWriter, r *http.Request) {
//	ctx := apm.DetachedContext(r.Context())
//	logger := GetLoggerWithContext(ctx).WithField("METHOD", "delete course according to request")
//	logger.Infof("API Delete course")
//
//	role, ok := r.Context().Value("role").(string)
//	if !ok {
//		http.Error(w, "Unable to get role/userName from token", http.StatusUnauthorized)
//		return
//	}
//
//	if role == "user" {
//		response := map[string]interface{}{
//			"message": "You are not allowed to this function",
//		}
//		w.Header().Set("Content-Type", "application/json")
//		w.WriteHeader(http.StatusOK)
//		json.NewEncoder(w).Encode(response)
//	}
//
//	keys := r.URL.Query()
//	userId := keys.Get("user_id")
//	if userId == "" {
//		// courseId is missing, return an error
//		log.Error("UserId parameter is missing")
//		http.Error(w, "UserId parameter is required ", http.StatusBadRequest)
//		return
//	}
//
//	err := userService.DeleteUserByIdService(userId, ctx)
//	if err != nil {
//		log.WithError(err).Errorf("Error delete user %s", userId)
//		http.Error(w, commonconstant.ErrUserNotExist, http.StatusBadRequest)
//	}
//}

//func handleSendDailyEmail(w http.ResponseWriter, r *http.Request) {
//	ctx := apm.DetachedContext(r.Context())
//	logger := GetLoggerWithContext(ctx).WithField("METHOD", "delete course according to request")
//	logger.Infof("API Delete course")
//
//	user.SendDailyEmail()
//}

func getStudentsByCourse(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD", "get students by course Id")
	logger.Infof("get students by course id")

	keys := r.URL.Query()
	courseId := keys.Get("courseId")

	studentList, err := userService.GetStudentByCourseId(courseId, ctx)
	if err != nil {
		log.WithError(err).Errorf("Unable to get students via course ID service")
		http.Error(w, "Unable to get students via course ID service", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Successful get student list for course",
		"data":    studentList,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handlePostStudentAttendance(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD POST", "post student attendance")
	logger.Infof("this API is used to add student attendance to database")

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Warningf("Error when reading from request")
		http.Error(w, "Invalid format", 252001)
		return
	}

	json.NewDecoder(r.Body)
	defer r.Body.Close()

	var rq api_request.StudentAttendanceRequest
	err = json.Unmarshal(bodyBytes, &rq)
	if err != nil {
		log.WithError(err).Errorf("Error marshaling body in getting student attendance request")
		http.Error(w, "Status internal Request", http.StatusInternalServerError)
		return
	}

	err = userService.AddStudentAttendanceService(rq, ctx)
	if err != nil {
		log.WithError(err).Errorf("Error service add student attendance")
		http.Error(w, "Unable to insert student attendance", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Successful add student attendance to database",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func updateStudentAttendanceStatus(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD PUT", "update student attendance")
	logger.Infof("this API is used to update student attendance to database")

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Warningf("Error when reading from request")
		http.Error(w, "Invalid format", 252001)
		return
	}

	json.NewDecoder(r.Body)
	defer r.Body.Close()

	var rq api_request.StudentAttendanceRequest
	err = json.Unmarshal(bodyBytes, &rq)

	err = userService.UpdateStudentAttendanceService(rq, ctx)
	if err != nil {
		log.WithError(err).Errorf("Error update student attendance request service")
		http.Error(w, "Error internal Request", http.StatusInternalServerError)
		return
	}

	log.Infof("Successful update student %s attendance", rq.StudentId)
	response := map[string]interface{}{
		"message": "Successful update students attendance to database",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleDeleteSessionByClassIs(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD DElETE", "Delete session by classId")
	logger.Infof("this API is used to delete class via classId")

	role, ok := r.Context().Value("role").(string)
	if !ok {
		http.Error(w, "Unable to get role/userName from token", http.StatusUnauthorized)
		return
	}

	if role == "user" {
		response := map[string]interface{}{
			"message": "You are not allowed to this function",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(response)
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Warningf("Error when reading from request")
		http.Error(w, "Invalid format", 252001)
		return
	}

	json.NewDecoder(r.Body)
	defer r.Body.Close()

	var rq api_request.DeleteClassInfo
	err = json.Unmarshal(bodyBytes, &rq)
	if err != nil {
		log.WithError(err).Errorf("Error marshal request to delete class")
		http.Error(w, "Error internal Request", http.StatusInternalServerError)
		return
	}

	err = classService.DeleteClassByClassId(rq, ctx)
	if err != nil {
		log.WithError(err).Errorf("Error delete class service")
		http.Error(w, "Error internal Request", http.StatusInternalServerError)
		return
	}

	log.Infof("Successful delete class %s", rq.ClassId)
	response := map[string]interface{}{
		"message": "Successful delete class via classId",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleCheckInAttendanceClass(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD POST", "Check in worker's attendance to calculate salary")
	logger.Infof("this API is used to delete class via classId")

	//role, ok := r.Context().Value("role").(string)
	//if !ok {
	//	http.Error(w, "Unable to get role/userName from token", http.StatusUnauthorized)
	//	return
	//}
	//
	//if role == "user" {
	//	response := map[string]interface{}{
	//		"message": "You are not allowed to this function",
	//	}
	//	w.Header().Set("Content-Type", "application/json")
	//	w.WriteHeader(http.StatusOK)
	//	json.NewEncoder(w).Encode(response)
	//}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Warningf("Error when reading from request")
		http.Error(w, "Invalid format", 252001)
		return
	}

	json.NewDecoder(r.Body)
	defer r.Body.Close()

	var rq api_request.CheckInAttendanceWorkerRequest
	err = json.Unmarshal(bodyBytes, &rq)
	if err != nil {
		log.WithError(err).Warningf("Error marshalling body from request update attendance worker")
		http.Error(w, "Invalid format", 252001)
		return
	}

	err = userService.CheckInWorkerAttendanceService(rq, rq.UserId, ctx)
	if err != nil {
		log.WithError(err).Warningf("Error check in attendance course to update attendance worker")
		http.Error(w, "Error internal", 252001)
		return
	}

	response := map[string]interface{}{
		"message": "Successful update worker attendance to database",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleDeleteStudent(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD DELETE", "handle delete student in course")
	logger.Infof("this API is used to delete class via classId")

	role, ok := r.Context().Value("role").(string)
	if !ok {
		http.Error(w, "Unable to get role/userName from token", http.StatusUnauthorized)
		return
	}

	if role == "user" {
		response := map[string]interface{}{
			"message": "You are not allowed to this function",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(response)
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Warningf("Error when reading from request")
		http.Error(w, "Invalid format", 252001)
		return
	}

	json.NewDecoder(r.Body)
	defer r.Body.Close()

	var rq api_request.DeleteStudentRequest
	err = json.Unmarshal(bodyBytes, &rq)
	if err != nil {
		log.WithError(err).Warningf("Error marshalling body from request update attendance worker")
		http.Error(w, "Invalid format", 252001)
		return
	}

	err = userService.DeleteStudentService(rq, ctx)
	if err != nil {
		log.WithError(err).Warningf("Error delete student service %s", err)
		http.Error(w, "Invalid format", 252001)
		return
	}

	response := map[string]interface{}{
		"message": "Successful delete student in course",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleFixStudentInformation(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD PUT", "handle modify student information")
	logger.Infof("this API is used to modify student information")

	role, ok := r.Context().Value("role").(string)
	if !ok {
		http.Error(w, "Unable to get role/userName from token", http.StatusUnauthorized)
		return
	}

	if role == "user" {
		response := map[string]interface{}{
			"message": "You are not allowed to this function",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(response)
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Warningf("Error when reading from request")
		http.Error(w, "Invalid format", 252001)
		return
	}

	json.NewDecoder(r.Body)
	defer r.Body.Close()

	var rq api_request.ModifyStudentRequest
	err = json.Unmarshal(bodyBytes, &rq)
	if err != nil {
		log.WithError(err).Warningf("Error marshalling body from request update attendance worker")
		http.Error(w, "Invalid format", 252001)
		return
	}

	err = userService.ModifyStudentInformation(rq, ctx)
	if err != nil {
		log.WithError(err).Warningf("Error update student information service %s", err)
		http.Error(w, "Invalid format", 252001)
		return
	}

	response := map[string]interface{}{
		"message": "Successful delete student in course",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleAddUserByPosition(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD", "handleRetryUserAccount")
	logger.Infof("Handle user account")

	role, ok := r.Context().Value("role").(string)
	if !ok {
		http.Error(w, "Unable to get role/userName from token", http.StatusUnauthorized)
		return
	}

	if role == "user" {
		response := map[string]interface{}{
			"message": "You are not allowed to this function",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(response)
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Warningf("Error when reading from request")
		http.Error(w, "Invalid format", 252001)
		return
	}

	json.NewDecoder(r.Body)
	defer r.Body.Close()

	var newUserRequest api_request.NewUserAddedByAdmin
	err = json.Unmarshal(bodyBytes, &newUserRequest)
	if err != nil {
		log.WithError(err).Warningf("Error when unmarshaling data from request")
		http.Error(w, "Status internal Request", http.StatusInternalServerError) // Return a internal server error
		return
	}

	// check whether the user exists in database
	userExistence, err := userService.GetByUserName(newUserRequest.UserName, newUserRequest.Email, "", ctx)
	if err != nil {
		log.WithError(err).Warningf("Error when get user data by username")
		http.Error(w, "Status internal Request", http.StatusInternalServerError)
		return
	}

	if userExistence.UserName == newUserRequest.UserName || userExistence.Email == newUserRequest.Email {
		log.Infof("UserName/Email existed")

		// Người dùng đã tồn tại, trả về một thông báo JSON cho phía frontend
		response := map[string]string{
			"message": "User/Email already exists",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict) // HTTP 409 Conflict status code
		json.NewEncoder(w).Encode(response)
		return
	}

	err = userService.InsertNewUserByJobPosition(newUserRequest, ctx)
	if err != nil {
		log.WithError(err).Errorf("Error insert new user")
		http.Error(w, "Status internal request", http.StatusInternalServerError)
		return
	}
}
