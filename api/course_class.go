package api

//func handleGetClassInformation(w http.ResponseWriter, r *http.Request) {
//	ctx := apm.DetachedContext(r.Context())
//	logger := GetLoggerWithContext(ctx).WithField("METHOD", "handleRetryUserAccount")
//	logger.Infof("Handle user account")
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
//	var classRequest batman.ClassRequest
//	err = json.Unmarshal(bodyBytes, &classRequest)
//	if err != nil {
//		log.WithError(err).Warningf("Error when unmarshaling data from request")
//		http.Error(w, "Status bad Request", http.StatusBadRequest) // Return a 400 Bad Request error
//		return
//	}
//
//	log.Infof("Get class information by class name")
//	classInfo, err := classService.GetClassInformationByClassName(classRequest, ctx)
//	if err != nil {
//		log.WithError(err).Errorf("Unable to get class information data")
//		return
//	}
//
//}
