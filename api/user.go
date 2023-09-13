package api

import (
	"database/sql"
	education_website "education-website"
	"education-website/client"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm"
	"io/ioutil"
	"log"
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
		logrus.WithError(err).Warningf("Error when reading from request")
		http.Error(w, "Invalid format", 252001)
		return
	}

	json.NewDecoder(r.Body)
	defer r.Body.Close()
	var userRequest client.UserRequest
	err = json.Unmarshal(bodyBytes, &userRequest)
	if err != nil {
		logrus.WithError(err).Warningf("Error when unmarshaling data from request")
		http.Error(w, "Status bad Request", http.StatusBadRequest) // Return a 400 Bad Request error
		return
	}

	user, err := userService.GetByUserName(userRequest.UserName, ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.WithError(err).Errorf("Username %s does not exist", user.UserName)
			return
		}
		logrus.WithError(err).Errorf("Error getting user info: %s", user.UserName)
		return
	}

	log.Printf("Successful get the user data")
	respondWithJSON(w, http.StatusOK, education_website.CommonResponse{Status: "SUCCESS", Descrition: "Success getting the user data"})
}
