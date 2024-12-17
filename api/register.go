package api

import (
	api_request "education-website/api/request"
	api_response "education-website/api/response"
	"education-website/entity/user"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.elastic.co/apm"
	"io/ioutil"
	"net/http"
	"time"
)

func handleTest(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD", "handleTest")
	logger.Infof("Handle test")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Hello world")
}

func handlerRegisterUser(w http.ResponseWriter, r *http.Request) {
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

	var registerRequest api_request.RegisterRequest
	err = json.Unmarshal(bodyBytes, &registerRequest)
	if err != nil {
		log.WithError(err).Warningf("Error when unmarshaling data from request")
		http.Error(w, "Status internal Request", http.StatusInternalServerError) // Return a internal server error
		return
	}

	// check whether the user exists in database
	userExistence, err := userService.GetByUserName(registerRequest.UserName, registerRequest.Email, "", ctx)
	if err != nil {
		log.WithError(err).Warningf("Error when get user data by username")
		http.Error(w, "Status internal Request", http.StatusInternalServerError)
		return
	}

	if userExistence.UserName == registerRequest.UserName || userExistence.Email == registerRequest.Email {
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

	userId, err := userService.InsertNewUser(registerRequest, ctx)
	if err != nil {
		log.WithError(err).Errorf("Error insert new user")
		http.Error(w, "Status internal request", http.StatusInternalServerError)
		return
	}

	generatedToken := jwtService.GenerateToken(user.UserEntity{
		UserId:   userId,
		UserName: registerRequest.UserName,
		FullName: registerRequest.FullName,
		Role:     "user",
	})

	userInfoWithToken := map[string]interface{}{
		"user": api_response.UserDto{
			UserId:       userId,
			UserName:     registerRequest.UserName,
			Email:        registerRequest.Email,
			Role:         "user",
			DOB:          registerRequest.DOB,
			JobPosition:  "Undefined",
			StartingDate: time.Now().Format("2006-01-02"),
		},
		"token": generatedToken,
	}

	response := map[string]interface{}{
		"message": "Đăng nhập thành công",
		"data":    userInfoWithToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleChangePassword(w http.ResponseWriter, r *http.Request) {
	userName, ok := r.Context().Value("username").(string)
	if !ok {
		http.Error(w, "Unable to get userName from token", http.StatusUnauthorized)
		return
	}

	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD", "handle changing password")
	logger.Infof("Handle user account")

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Warningf("Error when reading from request")
		http.Error(w, "Invalid format", 252001)
		return
	}

	json.NewDecoder(r.Body)
	defer r.Body.Close()

	var changePasswordRequest api_request.ChangePasswordRequest
	err = json.Unmarshal(bodyBytes, &changePasswordRequest)
	if err != nil {
		log.WithError(err).Warningf("Error when unmarshaling data from request")
		http.Error(w, "Status internal Request", http.StatusInternalServerError) // Return a internal server error
		return
	}

	err = userService.ChangePassword(changePasswordRequest, userName, ctx)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to change password: %s", err.Error())
		log.WithError(err).Errorf(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Update new password successful")
}
