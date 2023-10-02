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
	user, err := userService.GetByUserName(userRequest.UserName, userRequest.Email, ctx)
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
	req := keys.Get("user") // get user request, all or one

	role, ok := r.Context().Value("role").(string)
	if !ok {
		http.Error(w, "Unable to get role from token", http.StatusUnauthorized)
		return
	}

	if role == "user" {
		userSalaryReport, err := userService.
	}

	if req == "all" {
		if role == "user" {
			response := map[string]string{
				"message": "You are not allowed to access to this function",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden) // HTTP 409 Conflict status code
			json.NewEncoder(w).Encode(response)
			return
		}
	} else {
		// user request access to his/her own salary table
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.WithError(err).Warningf("Error when reading from request")
			http.Error(w, "Invalid format", 252001)
			return
		}

		json.NewDecoder(r.Body)
		defer r.Body.Close()

		var registerRequest api_request.RegisterRequest
		err = json.Unmarshal(bodyBytes, &registerRequest)
	}

}
