package api

import (
	batman "education-website"
	api_request "education-website/api/request"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go.elastic.co/apm"
	"io/ioutil"
	"net/http"
)

func handlerUserAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}

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
	role, ok := r.Context().Value("role").(string)
	userName, ok := r.Context().Value("username").(string)
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
	var salaryRequest api_request.SalaryRequest
	err = json.Unmarshal(bodyBytes, &salaryRequest)
	if err != nil {
		log.WithError(err).Warningf("Error when unmarshaling data from request")
		http.Error(w, "Status bad Request", http.StatusBadRequest) // Return a 400 Bad Request error
		return
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

	keys := r.URL.Query() // this variable is used for searching engine for leader
	role, ok := r.Context().Value("role").(string)
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
	var newSalaryInfo api_request.ModifySalaryConfRequest
	err = json.Unmarshal(bodyBytes, &newSalaryInfo)
	if err != nil {
		log.WithError(err).Warningf("Error when unmarshaling data from request")
		http.Error(w, "Status bad Request", http.StatusBadRequest) // Return a 400 Bad Request error
		return
	}

	if role == "user" {
		response := map[string]interface{}{
			"message": "You are not allowed to access to this function",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
	} else {

	}
}
