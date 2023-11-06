package api

import (
	batman "education-website"
	api_request "education-website/api/request"
	api_response "education-website/api/response"
	"education-website/service/user"
	"encoding/json"
	"fmt"
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
			userSalaryReport, err := userService.GetSalaryInformation("", salaryRequest.Month, salaryRequest.Year, ctx)
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

	err := r.ParseMultipartForm(10 << 20) // 10 MB limit for the uploaded file
	if err != nil {
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

	err = userService.ImportStudentsByExcel(file, ctx)
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
