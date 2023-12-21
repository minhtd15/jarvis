package api

import (
	api_request "education-website/api/request"
	api_response "education-website/api/response"
	"education-website/entity/user"
	user2 "education-website/service/user"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go.elastic.co/apm"
	"io/ioutil"
	"net/http"
)

func handlerLoginUser(w http.ResponseWriter, r *http.Request) {
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

	// this is the information that the user type in front end
	var userRequest api_request.LoginRequest
	err = json.Unmarshal(bodyBytes, &userRequest)
	if err != nil {
		log.WithError(err).Warningf("Error when unmarshaling data from request")
		http.Error(w, "Status internal Request", http.StatusInternalServerError) // Return a internal server error
		return
	}

	// get the user from request from database
	userEntityInfo, err := userService.GetUserNamePassword(userRequest, ctx)
	if err != nil {
		log.WithError(err).Warningf("Error verify Username and Password for user")
		http.Error(w, "Status internal Request", http.StatusInternalServerError) // Return a internal server error
		return
	}

	// verify the user password
	checkPasswordSimilarity, err := authService.VerifyUser(userRequest, *userEntityInfo)
	if err != nil {
		log.WithError(err).Errorf("Invalid username or password for user: %s", userEntityInfo.UserName)
		http.Error(w, "wrong password/username", http.StatusBadRequest)
		return
	}

	// if the username exist and the password is true, generate token to send back to backend
	if checkPasswordSimilarity != nil {
		generatedToken := jwtService.GenerateToken(user.UserEntity{
			UserId:   userEntityInfo.UserId,
			UserName: userEntityInfo.UserName,
			FullName: userEntityInfo.FullName,
			Role:     userEntityInfo.Role,
		})

		// Gán token vào đối tượng model.User
		userInfoWithToken := map[string]interface{}{
			"user": api_response.UserDto{
				UserId:       userEntityInfo.UserId,
				UserName:     userEntityInfo.UserName,
				FullName:     userEntityInfo.FullName,
				Email:        userEntityInfo.Email,
				Role:         userEntityInfo.Role,
				DOB:          userEntityInfo.DOB,
				JobPosition:  userEntityInfo.JobPosition,
				StartingDate: userEntityInfo.StartDate,
			},
			"token": generatedToken,
		}

		response := map[string]interface{}{
			"message": "Đăng nhập thành công",
			"data":    userInfoWithToken,
		}

		// Trả về thông tin người dùng cùng với token
		// Trả về phản hồi JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	} else {
		// return cannot find user
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid password or username")
	}
}

func handleSendEmailForgotPassword(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD", "handle forgot password")
	logger.Infof("Handle forgot password")

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Warningf("Error when reading from request")
		http.Error(w, "Invalid format", 252001)
		return
	}

	json.NewDecoder(r.Body)
	defer r.Body.Close()

	// this is the information that the user type in front end
	var userRequest api_request.ForgotPasswordRequest
	err = json.Unmarshal(bodyBytes, &userRequest)
	if err != nil {
		log.WithError(err).Warningf("Error when unmarshaling data from request")
		http.Error(w, "Status internal Request", http.StatusInternalServerError) // Return a internal server error
		return
	}

	email := userRequest.Email

	checkEmailExistence, err := userService.CheckEmailExistenceService(email, ctx)
	if err != nil || !checkEmailExistence {
		log.WithError(err).Errorf("Email probably not existed %s", email)
		http.Error(w, "Email not existed, please create new account", http.StatusBadRequest)
		return
	}

	// update a digit code to db
	digitCode, err := userService.PostNewForgotPasswordCode(email, ctx)
	if err != nil {
		log.WithError(err).Errorf("Cannot update digi code to db", email)
		http.Error(w, "Cannot update digi code to db", http.StatusBadRequest)
		return
	}

	// send the digit code to the email
	user2.SendDailyEmail(email, *digitCode)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Success send digit code to email")

}

func handleCheckDigitCodeForgotPassword(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD", "handle forgot password")
	logger.Infof("Handle forgot password to check digit code")

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Warningf("Error when reading from request")
		http.Error(w, "Invalid format", 252001)
		return
	}

	json.NewDecoder(r.Body)
	defer r.Body.Close()

	// this is the information that the user type in front end
	var userRequest api_request.CheckDigitCode
	err = json.Unmarshal(bodyBytes, &userRequest)
	if err != nil {
		log.WithError(err).Warningf("Error when unmarshaling data from request")
		http.Error(w, "Status internal Request", http.StatusInternalServerError) // Return a internal server error
		return
	}

	code := userRequest.DigitCode
	email := userRequest.Email

	check, err := userService.CheckFitDigitCode(email, code, ctx)
	if err != nil {
		log.WithError(err).Errorf("unable to check on db digit code")
		http.Error(w, "Unable to check on db digit code", http.StatusBadRequest)
		return
	}

	if !*check {
		log.WithError(err).Errorf("Wrong code")
		http.Error(w, "Wrong code", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Success send verify digit code")
}

func handleSetNewPassword(w http.ResponseWriter, r *http.Request) {
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

	var newPasswordRequest api_request.NewPasswordRequest
	err = json.Unmarshal(bodyBytes, &newPasswordRequest)
	if err != nil {
		log.WithError(err).Warningf("Error when unmarshaling data from request")
		http.Error(w, "Status internal Request", http.StatusInternalServerError) // Return a internal server error
		return
	}

	userInfo, err := userService.UpdateNewPasswordInfo(newPasswordRequest.NewPassword, newPasswordRequest.Email, ctx)
	if err != nil {
		log.WithError(err).Warningf("Error when insert new password data from request")
		http.Error(w, "Status internal Request", http.StatusInternalServerError) // Return a internal server error
		return
	}

	generatedToken := jwtService.GenerateToken(user.UserEntity{
		UserId:       userInfo.UserId,
		UserName:     userInfo.UserName,
		FullName:     userInfo.FullName,
		Role:         userInfo.Role,
		DOB:          userInfo.DOB,
		JobPosition:  userInfo.JobPosition,
		StartingDate: userInfo.StartingDate,
	})

	userInfoWithToken := map[string]interface{}{
		"user": api_response.UserDto{
			UserId:       userInfo.UserId,
			UserName:     userInfo.UserName,
			Email:        newPasswordRequest.Email,
			Role:         userInfo.Role,
			DOB:          userInfo.DOB,
			JobPosition:  userInfo.JobPosition,
			StartingDate: userInfo.StartingDate,
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
