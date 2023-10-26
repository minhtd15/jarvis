package api

import (
	api_request "education-website/api/request"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go.elastic.co/apm"
	"io/ioutil"
	"net/http"
)

func handleInsertNewClass(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD", "handle get salary information")
	logger.Infof("Handle user account")

	//role, ok := r.Context().Value("role").(string)
	//userName, ok := r.Context().Value("username").(string)
	//if !ok {
	//	http.Error(w, "Unable to get role/userName from token", http.StatusUnauthorized)
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
	var salaryRequest api_request.SalaryRequest
	err = json.Unmarshal(bodyBytes, &salaryRequest)
	if err != nil {
		log.WithError(err).Warningf("Error when unmarshaling data from request")
		http.Error(w, "Status bad Request", http.StatusBadRequest) // Return a 400 Bad Request error
		return
	}
}

func handleGetClassInformation(w http.ResponseWriter, r *http.Request) {
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
	var classInfoRequest api_request.CourseInfoRequest
	err = json.Unmarshal(bodyBytes, &classInfoRequest)
	if err != nil {
		log.WithError(err).Warningf("Unable to unmarshal data from request")
		http.Error(w, "Status bad Request", http.StatusBadRequest) // Return a 400 Bad Request error
		return
	}

	classInfoResponse, err := classService.GetCourseInformationByClassName(classInfoRequest, ctx)
	if err != nil {
		log.WithError(err).Warningf("Unable to unmarshal data from request")
		http.Error(w, "Status bad Request", http.StatusBadRequest) // Return a 400 Bad Request error
		return
	}

	response := map[string]interface{}{
		"message": "Successful getting class information: " + classInfoRequest.CourseId,
		"data":    classInfoResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
