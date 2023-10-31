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
	logger := GetLoggerWithContext(ctx).WithField("METHOD", "handle insert new course information")
	logger.Infof("Handle new course information")

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
	var newCourseRequest api_request.NewCourseRequest
	err = json.Unmarshal(bodyBytes, &newCourseRequest)
	if err != nil {
		log.WithError(err).Warningf("Error when unmarshaling data from request")
		http.Error(w, "Status bad Request", http.StatusInternalServerError)
		return
	}

	err = classService.AddNewClass(newCourseRequest, ctx)
	if err != nil {
		log.WithError(err).Errorf("Unable to create new class")
		http.Error(w, "Status bad Request", http.StatusInternalServerError)
		return
	}
}

func handleGetClassInformation(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD", "handle get course information")
	logger.Infof("Handle course information")

	keys := r.URL.Query()
	courseId := keys.Get("classId")

	classInfoRequest := api_request.CourseInfoRequest{
		CourseId: courseId,
	}

	classInfoResponse, err := classService.GetCourseInformationByClassName(classInfoRequest, ctx)
	if err != nil {
		log.WithError(err).Warningf("Unable to unmarshal data from request")
		http.Error(w, "Status bad Request", http.StatusInternalServerError) // Return a 400 Bad Request error
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

func handleGetAllCourseInformation(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD", "handle get all courses information")
	logger.Infof("Get all courses information")

	allCoursesInfo, err := classService.GetAllCourses(ctx)
	if err != nil {
		log.WithError(err).Errorf("Error getting all courses information")
		http.Error(w, "Status bad Request", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Successful getting all course",
		"data":    allCoursesInfo,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
