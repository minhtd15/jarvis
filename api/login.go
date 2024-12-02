package api

import (
	"bytes"
	api_request "education-website/api/request"
	api_response "education-website/api/response"
	"education-website/entity/user"
	"encoding/base64"
	"encoding/json"
	"fmt"
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

func handlerLogin(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD", "handleLogin")
	logger.Infof("Handle login WSO2")

	// Lấy mã xác thực từ truy vấn
	// Lấy thông tin đăng nhập từ yêu cầu HTTP
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Kiểm tra xem người dùng đã gửi đủ thông tin đăng nhập chưa
	if username == "" || password == "" {
		http.Error(w, "Missing username or password", http.StatusBadRequest)
		return
	}

	// Cấu hình Client ID và Client Secret từ WSO2
	clientID := "<YourClientID>"
	clientSecret := "<YourClientSecret>"

	// Tạo Basic Auth Header từ ClientID và ClientSecret
	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	// Tạo dữ liệu yêu cầu POST để lấy token
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", username)
	data.Set("password", password)
	data.Set("scope", "openid")

	// Gửi yêu cầu POST tới WSO2 để lấy token
	tokenURL := "https://<API_GATEWAY_URL>/oauth2/token"
	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		logger.Errorf("Error creating request: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Đặt header cho yêu cầu
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("Error sending request to WSO2: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Đọc phản hồi từ WSO2
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Error reading response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Kiểm tra mã trạng thái của phản hồi
	if resp.StatusCode != http.StatusOK {
		logger.Errorf("Failed to get token: %s", body)
		http.Error(w, "Failed to get token from WSO2", http.StatusUnauthorized)
		return
	}

	// Hiển thị token cho người dùng (bạn có thể lưu trữ token và sử dụng nó cho các API khác)
	logger.Infof("Successfully obtained token: %s", body)
	w.Write(body) // Gửi token về cho người dùng
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

	clientID := "Bw0g06uSGL317MdJHyjSsU39930a"      // Consumer Key từ WSO2 Developer Portal
	redirectURI := "http://localhost:4200/callback" // URI mà người dùng sẽ được chuyển hướng sau khi xác thực thành công
	scope := "openid"                               // Các quyền mà bạn yêu cầu
	authURL := fmt.Sprintf("https://localhost:9443/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=%s",
		clientID, redirectURI, scope, state)

	//auth0Domain := "dev-5wpln5bbc476iydk.us.auth0.com"
	//clientID := "wDIUqMHxB4XxNDzOUbhtDO66Qd72kJ8a"
	//redirectURI := "http://localhost:8081/e/v1/callback"
	//redirectURI := "http://localhost:3031/crm-tiw/loading"
	//loginURL := "https://" + auth0Domain + "/authorize" +
	//	"?response_type=code" +
	//	"&client_id=" + clientID +
	//	"&redirect_uri=" + redirectURI +
	//	"&scope=openid profile email" +
	//	"&state=" + state

	// Thực hiện chuyển hướng
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")                   // Cho phép từ miền của bạn
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // Các phương thức được phép
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")       // Các header cần thiết
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	// Phần còn lại của mã
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
	var requestData struct {
		Code string `json:"code"`
	}

	// Giải mã dữ liệu JSON từ body
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		log.WithError(err).Error("Failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	authorizationCode := requestData.Code
	if authorizationCode == "" {
		log.Error("Authorization code is missing")
		http.Error(w, "Authorization code is missing", http.StatusBadRequest)
		return
	}

	log.Infof("Authorization code: %s", authorizationCode)

	// Gọi hàm để trao đổi mã xác thực với Auth0 lấy token
	response, err := exchangeAuthCode(authorizationCode)
	if err != nil {
		log.WithError(err).Errorf("Failed to exchange authorization code for tokens: %s", err)
		http.Error(w, "Failed to exchange authorization code for tokens", http.StatusInternalServerError)
		return
	}
	log.Infof("Response: %v", response)

	// Trả về phản hồi JSON
	log.Info("Successfully logged in via Auth0")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func exchangeAuthCode(code string) (*string, error) {
	// Cấu trúc yêu cầu trao đổi mã xác thực
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", "Bw0g06uSGL317MdJHyjSsU39930a")
	data.Set("client_secret", "EPzKR3eoLr1HfPyir1lCAi2AaPsa")
	data.Set("code", code)
	data.Set("redirect_uri", "http://localhost:4200/callback")

	// Gửi yêu cầu POST đến WSO2 để trao đổi mã xác thực

	resp, err := http.PostForm("https://localhost:9443/oauth2/token", data)
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
		log.WithError(err).Errorf("Failed to parse ID token: %s", err)
		return nil, err
	}
	nickname := ""
	if value, ok := claims["nickname"].(string); ok {
		nickname = value
	}
	log.Infof("Nickname: %s", nickname)

	return nil, nil
}
