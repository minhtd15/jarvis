package api

import (
	api_request "education-website/api/request"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.elastic.co/apm"
	"io/ioutil"
	"net/http"
	"strconv"
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
		http.Error(w, "Status internal server error", http.StatusInternalServerError)
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

func handleClassFromToDateById(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD", "handle get default class information")
	logger.Infof("Get default class information")

	userId, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unable to get userId from token", http.StatusUnauthorized)
		return
	}

	keys := r.URL.Query()
	fromDate := keys.Get("fromDate")
	toDate := keys.Get("toDate")

	courseType, err := classService.GetCourseType(ctx)
	if err != nil {
		log.WithError(err).Errorf("Unable to get course type")
		http.Error(w, "unable to get course type api", http.StatusInternalServerError)
		return
	}

	schedule, err := classService.GetFromToSchedule(fromDate, toDate, userId, courseType, ctx)
	if err != nil {
		log.WithError(err).Errorf("Error getting all schedule api for user: %s", userId)
		http.Error(w, "Unable to get schedule", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Successful getting schedule information",
		"data":    schedule,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

//func handleCheckInAttendanceClass(w http.ResponseWriter, r *http.Request) {
//	ctx := apm.DetachedContext(r.Context())
//	logger := GetLoggerWithContext(ctx).WithField("METHOD", "handle checkin class")
//	logger.Infof("Start to check in class API")
//
//	userTokenId, ok := r.Context().Value("user_id").(string)
//	if !ok {
//		http.Error(w, "Unable to get userId from token", http.StatusUnauthorized)
//		return
//	}
//
//	bodyBytes, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		log.WithError(err).Warningf("Error when reading from request")
//		http.Error(w, "Invalid format", 252001)
//		return
//	}
//
//	json.NewDecoder(r.Body)
//	defer r.Body.Close()
//	var checkInAttendance api_request.ChecInAttendanceClassRequest
//	err = json.Unmarshal(bodyBytes, &checkInAttendance)
//	if err != nil {
//		log.WithError(err).Errorf("Error marshal Check In Attendance Request: %s", err)
//		http.Error(w, "Error marshalling checkin attendance request", http.StatusInternalServerError)
//		return
//	}
//
//	// get the start time of the course to check whether
//	classInformation, err := classService.GetClassInformationByClassId(checkInAttendance.ClassId, ctx)
//	if err != nil {
//		log.WithError(err).Errorf("Error getting class information from db")
//		http.Error(w, "Error getting class information from db", http.StatusInternalServerError)
//		return
//	}
//
//	// add information to ATTENDANCE_HISTORY table
//
//	// return success check in attendance to user
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//}

func checkInStudentAttendance(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD", "handle get student attendance by course Id ")
	logger.Infof("Get class all sessions")

	keys := r.URL.Query()
	courseId := keys.Get("courseId")

	result, err := userService.GetCourseSessionsService(courseId, ctx)
	if err != nil {
		log.WithError(err).Errorf("Unable to get course sessions service")
		http.Error(w, "unable to get course session api", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Successful getting schedule information",
		"data":    result,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleGetMySchedule(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD", "handle user every class that this person in charge")
	logger.Infof("Get course in charge")

	username, ok := r.Context().Value("username").(string)
	if !ok {
		http.Error(w, "Unable to get userId from token", http.StatusUnauthorized)
		return
	}

	result, err := userService.GetAllInChargeCourse(username, ctx)
	if err != nil {
		log.WithError(err).Errorf("Unable to get course in charge service")
		http.Error(w, "unable to get course that this person in charge api", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Successful getting user course in charge",
		"data":    result,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleGetAllSessionsByCourseId(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD", "handle get course all sessions, students, ")
	logger.Infof("Get class all sessions")

	keys := r.URL.Query()
	courseId := keys.Get("courseId")

	rs, err := classService.GetAllSessionsByCourseIdService(courseId, ctx)
	if err != nil {
		log.WithError(err).Errorf("Unable to get sessions by course Id Service")
		http.Error(w, "Unable to get sessions by course Id", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":   "Successful getting sessions for course",
		"course_id": courseId,
		"data":      rs,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleFixCourseInformation(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD PUT", "fix course information")
	logger.Infof("this API is used to modify course information")

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

	var rq api_request.ModifyCourseInformation
	err = json.Unmarshal(bodyBytes, &rq)
	if err != nil {
		log.WithError(err).Errorf("Error marshal request to fix course information")
		http.Error(w, "Error marshal request to fix course information", http.StatusInternalServerError)
		return
	}

	err = classService.FixCourseInformationService(rq, ctx)
	if err != nil {
		log.WithError(err).Errorf("Error in service modify course information")
		http.Error(w, "Error in service modify course information", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Successful fix course information",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func AddNoteByClassId(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD POST", "add note to class by classId")
	logger.Infof("this API is used to add note to class by classId")

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Warningf("Error when reading from request")
		http.Error(w, "Invalid format", 252001)
		return
	}

	role, ok := r.Context().Value("role").(string)
	if !ok {
		http.Error(w, "Unable to get role/userName from token", http.StatusUnauthorized)
		return
	}

	json.NewDecoder(r.Body)
	defer r.Body.Close()
	var noteRequest api_request.AddNoteRequest
	err = json.Unmarshal(bodyBytes, &noteRequest)
	if err != nil {
		log.WithError(err).Warningf("Error when unmarshaling data from request")
		http.Error(w, "Status bad Request", http.StatusInternalServerError)
		return
	}

	intNumber, err := strconv.ParseInt(noteRequest.ClassId, 10, 64)
	if err != nil {
		// Handle the error if the conversion fails
		fmt.Println("Error converting string to int:", err)
		return
	}

	taList, err := classService.GetTAListService(int(intNumber), ctx)

	addTA := Difference(noteRequest.TaList, taList)    // get the new TA to add
	deleteTA := Difference(taList, noteRequest.TaList) // get the new TA to delete

	log.Infof("ADD TA: %v", addTA)
	log.Infof("DELETE TA: %v", deleteTA)

	if noteRequest.Check && role == "user" {
		response := map[string]interface{}{
			"message": "You are not allowed to access to this function",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(response)
	}

	err = classService.AddNoteService(noteRequest, addTA, deleteTA, ctx)
	if err != nil {
		log.WithError(err).Warningf("Error service add note")
		http.Error(w, "Error service add note", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Successful getting add note to class by classId",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func Difference(arr1, arr2 []string) []string {
	diffMap := make(map[string]bool)

	// Mark all elements in arr2 as true
	for _, num := range arr2 {
		diffMap[num] = true
	}

	// Filter elements in arr1 that are not in arr2
	var result []string
	for _, num := range arr1 {
		if !diffMap[num] {
			result = append(result, num)
		}
	}

	return result
}

func handleGetCheckInWorkerHistory(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD GET", "GET the check in history of course")
	logger.Infof("this API is used to GET the check in history of course")

	keys := r.URL.Query()
	courseId := keys.Get("courseId")

	rs, err := classService.GetCheckInHistoryByCourseId(courseId, ctx)
	if err != nil {
		log.WithError(err).Warningf("Error get check in history for course %s", courseId)
		http.Error(w, "Error get check in history", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Successful getting check in history for course",
		"data":    rs,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handlePostSubClass(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD POST", "add sub class")
	logger.Infof("this API is used to add sub class to course")

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Warningf("Error when reading from request")
		http.Error(w, "Invalid format", 252001)
		return
	}

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
		http.Error(w, "Status bad Request", http.StatusForbidden)
		json.NewEncoder(w).Encode(response)
		return
	}

	json.NewDecoder(r.Body)
	defer r.Body.Close()
	var subClassRq api_request.NewSubClassRequest
	err = json.Unmarshal(bodyBytes, &subClassRq)
	if err != nil {
		log.WithError(err).Warningf("Error when unmarshaling data from request")
		http.Error(w, "Status bad Request", http.StatusInternalServerError)
		return
	}

	err = classService.AddSubClassService(subClassRq, ctx)
	if err != nil {
		log.WithError(err).Warningf("Error add new sub class")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Successful add sub class",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleGetSubClass(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD GET", "get all sub class for course")
	logger.Infof("this API is used to add sub class to course")

	keys := r.URL.Query()
	courseId := keys.Get("courseId")

	rs, err := classService.GetSubClassByCourseId(courseId, ctx)
	if err != nil {
		log.WithError(err).Errorf("Unable to get sub class for course %s", courseId)
		http.Error(w, "Unable to get sub class by course Id", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Successful getting sub class for course",
		"data":    rs,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleDeleteSubClass(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD DELETE", "delete sub class")
	logger.Infof("this API is used to delete sub class to course")

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
	var deleteSubClassRq api_request.DeleteSubClassRequest
	err = json.Unmarshal(bodyBytes, &deleteSubClassRq)
	if err != nil {
		log.WithError(err).Warningf("Error when unmarshaling data from request")
		http.Error(w, "Status bad Request", http.StatusInternalServerError)
		return
	}

	err = classService.DeleteSubClassService(deleteSubClassRq, ctx)
	if err != nil {
		log.WithError(err).Warningf("Error delete sub class")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Successful delete sub class",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
