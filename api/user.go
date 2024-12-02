package api

import (
	batman "education-website"
	"encoding/json"
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
