package api

import (
	"context"
	api_request "education-website/api/request"
	api_response "education-website/api/response"
	"education-website/entity/user"
	user2 "education-website/service/user"
	"encoding/base64"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"go.elastic.co/apm"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
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

func handleGetStudentPaymentStatusByCourseId(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD GET", "get student payment status by course Id")
	logger.Infof("this API is used to get student payment status by course Id")

	keys := r.URL.Query()
	courseId := keys.Get("courseId")

	rs, err := userService.GetStudentPaymentStatusByCourseIdService(courseId, ctx)
	if err != nil {
		log.Errorf("Unable to get student payment status by course Id: %s; err : %v", courseId, err)
		http.Error(w, "Unable to get student payment status by course Id", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Successful getting student payment status by course Id",
		"data":    rs,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func loginViaThirdParty(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD", r.Method).WithField("API", "login via third party")
	logger.Infof("This API is used to login via third party: Auth0")

	state, err := generateState()
	if err != nil {
		// Xử lý lỗi khi tạo state
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	auth0Domain := "dev-5wpln5bbc476iydk.us.auth0.com"
	clientID := "wDIUqMHxB4XxNDzOUbhtDO66Qd72kJ8a"
	redirectURI := "http://localhost:8081/e/v1/callback"

	loginURL := "https://" + auth0Domain + "/authorize" +
		"?response_type=code" +
		"&client_id=" + clientID +
		"&redirect_uri=" + redirectURI +
		"&scope=openid profile email" +
		"&state=" + state

	// Thực hiện chuyển hướng
	http.Redirect(w, r, loginURL, http.StatusTemporaryRedirect)
}

func generateState() (string, error) {
	// Tạo một chuỗi state ngẫu nhiên
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	state := base64.StdEncoding.EncodeToString(randomBytes)
	return state, nil
}

func handleCallback(w http.ResponseWriter, r *http.Request) {

	// Lấy mã xác thực từ truy vấn
	authorizationCode := r.URL.Query().Get("code")
	//role := r.URL.Query().Get("role")

	// Thực hiện giao dịch mã xác thực để nhận mã thông báo từ Auth0
	response, err := exchangeAuthCode(authorizationCode)
	if err != nil {
		log.WithError(err).Errorf("Failed to exchange authorization code for tokens: %s", err)
		http.Error(w, "Failed to exchange authorization code for tokens", http.StatusInternalServerError)
		return
	}

	// Trả về phản hồi JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func exchangeAuthCode(code string) (map[string]interface{}, error) {
	// Cấu trúc yêu cầu trao đổi mã xác thực
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", "wDIUqMHxB4XxNDzOUbhtDO66Qd72kJ8a")
	data.Set("client_secret", "ts5cjzjMmPIvxytIjeBCeW88Re8HRdFL-A9w1PPQ2SxHyUiXVtHNsAacRGNOehd7")
	data.Set("code", code)
	data.Set("redirect_uri", "http://localhost:8081/e/v1/callback")

	// Gửi yêu cầu POST đến Auth0 để trao đổi mã xác thực
	resp, err := http.PostForm("https://dev-5wpln5bbc476iydk.us.auth0.com/oauth/token", data)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Đọc và phân tích phản hồi JSON
	var tokenResp struct {
		AccessToken string `json:"access_token"`
		IDToken     string `json:"id_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	// Giải mã IDToken để truy cập thông tin vai trò
	claims := jwt.MapClaims{}
	_, _, err = new(jwt.Parser).ParseUnverified(tokenResp.IDToken, &claims)
	if err != nil {
		return nil, err
	}
	nickname := ""
	if value, ok := claims["nickname"].(string); ok {
		nickname = value
	}
	log.Infof("Nickname: %s", nickname)

	// get role from token
	//role := getRoleByAuthManagementAPI()
	//log.Infof("Role: %s", role)

	// get user other information on database
	userInfo, err := userService.GetByNicknameService(nickname, context.Context(context.Background()))
	if err != nil {
		log.WithError(err).Errorf("Failed to get user information by nickname: %s", err)
		return nil, err
	}

	generatedToken := jwtService.GenerateToken(user.UserEntity{
		UserId:   userInfo.UserId,
		UserName: userInfo.UserName,
		FullName: userInfo.FullName,
		Role:     userInfo.Role,
	})

	userInfoWithToken := map[string]interface{}{
		"user": api_response.UserDto{
			UserId:      userInfo.UserId,
			UserName:    userInfo.UserName,
			FullName:    userInfo.FullName,
			Email:       userInfo.Email,
			Role:        userInfo.Role,
			DOB:         userInfo.DOB,
			JobPosition: userInfo.JobPosition,
		},
		"token": generatedToken,
	}

	return userInfoWithToken, nil
}

//func getRoleByAuthManagementAPI() string {
//	// Get these from your Auth0 Application Dashboard.
//	// The application needs to be a Machine To Machine authorized
//	// to request access tokens for the Auth0 Management API,
//	// with the desired permissions (scopes).
//	domain := "dev-5wpln5bbc476iydk.us.auth0.com"
//	clientID := "c3ySkmVVxNqrCx72z5eSKUXu039hA0Br"
//	clientSecret := "o7vRruLQ_BQWV9g5LEyfdkWfBYmX-BRFvhXL_tlbWxsX-g1-q2TyjOCAks8RU24W"
//
//	// Initialize a new client using a domain, client ID and client secret.
//	// Alternatively you can specify an access token:
//	// `management.WithStaticToken("token")`
//	auth0API, err := management.New(
//		domain,
//		management.WithClientCredentials(context.TODO(), clientID, clientSecret), // Replace with a Context that better suits your usage
//	)
//	if err != nil {
//		log.Fatalf("failed to initialize the auth0 management API client: %+v", err)
//	}
//
//	// Now we can interact with the Auth0 Management API.
//	// Example: Creating a new client.
//	client := &management.Client{
//		Name:        auth0.String("My Client"),
//		Description: auth0.String("Client created through the Go SDK"),
//	}
//
//	// The passed in client will get hydrated with the response.
//	// This means that after this request, we will have access
//	// to the client ID on the same client object.
//	err = auth0API.Client.Create(context.TODO(), client) // Replace with a Context that better suits your usage
//	if err != nil {
//		log.Fatalf("failed to create a new client: %+v", err)
//	}
//
//	// Make use of the getter functions to safely access
//	// fields without causing a panic due nil pointers.
//	log.Printf(
//		"Created an auth0 client successfully. The ID is: %q",
//		client.GetClientID(),
//	)
//	return client.GetClientID()
//}
