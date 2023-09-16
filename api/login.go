package api

import (
	api_request "education-website/api/request"
	api_response "education-website/api/response"
	"education-website/entity/user"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go.elastic.co/apm"
	"go.elastic.co/apm/model"
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
	var userRequest api_request.LoginRequest
	err = json.Unmarshal(bodyBytes, &userRequest)
	if err != nil {
		log.WithError(err).Warningf("Error when unmarshaling data from request")
		http.Error(w, "Status internal Request", http.StatusInternalServerError) // Return a internal server error
		return
	}

	// get the user from request from database
	userLoginInfo, err := userService.GetUserNamePassword(userRequest, ctx)
	if err != nil {
		log.WithError(err).Warningf("Error verify Username and Password for user")
		http.Error(w, "Status internal Request", http.StatusInternalServerError) // Return a internal server error
		return
	}

	// verify the user password
	checkPasswordSimilarity, err := authService.VerifyUser(userRequest, *userLoginInfo)
	if err != nil {
		log.Infof("Invalid username or password for user: %s", userLoginInfo.UserName)
		return
	}

	// if the username exist and the password is true, generate token to send back to backend
	if _, ok := checkPasswordSimilarity.(model.User); ok {
		generatedToken := jwtService.GenerateToken(user.UserEntity{
			UserName: userLoginInfo.UserName,
 			Role:     userLoginInfo.Role,
		})

		// Gán token vào đối tượng model.User
		userInfoWithToken := map[string]interface{}{
			"user": api_response.UserDto{
				UserId:       userLoginInfo.UserId,
				UserName:     userLoginInfo.UserName,
				Email:        userLoginInfo.Email,
				Role:         userLoginInfo.Role,
				DOB:          userLoginInfo.DOB,
				JobPosition:  userLoginInfo.JobPosition,
				StartingDate: userLoginInfo.StartDate,
			},
			"token": generatedToken,
		}

		// Trả về thông tin người dùng cùng với token
		// Trả về phản hồi JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(userInfoWithToken)
	}

	// return cannot find user

	chua viet dau, mai viet tiep
}
