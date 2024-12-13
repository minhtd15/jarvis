package api

import (
	api_request "education-website/api/request"
	"education-website/rabbitmq"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	_ "github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"go.elastic.co/apm"
	"io"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Cho phép kết nối từ bất kỳ origin nào
		return true
	},
}

func handleGetSports(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD GET", "get all sports")
	logger.Infof("this API is used to get all sports")

	rs, err := classService.GetSportsService(ctx)
	if err != nil {
		log.WithError(err).Errorf("Unable to get all sports")
		http.Error(w, "Unable to get all sports", http.StatusInternalServerError)
	}

	//response := map[string]interface{}{
	//	"message": "Successful getting all sports",
	//	"data":    rs,
	//}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rs)
}

func uploadImage(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD POST", "upload image")
	logger.Infof("this API is used to upload image")

	// Kiểm tra phương thức
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Lấy file từ body request
	file, _, err := r.FormFile("imageFile")
	if err != nil {
		logger.WithError(err).Warningf("Error retrieving file from request")
		http.Error(w, "Error retrieving file from request", http.StatusBadRequest)
		return
	}

	imgData, err := io.ReadAll(file)
	if err != nil {
		logger.WithError(err).Error("Error reading file data")
		http.Error(w, "Error reading file data", http.StatusInternalServerError)
		return
	}

	defer file.Close()

	// Đọc thông tin về hình ảnh từ form
	var rq api_request.UploadImageRequest
	if err := r.ParseMultipartForm(10 << 20); err != nil { // Giới hạn kích thước form
		logger.WithError(err).Warningf("Error parsing multipart form")
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Lấy sportId từ form
	sportIdStr := r.FormValue("sportId")
	if sportIdStr == "" {
		http.Error(w, "sportId is required", http.StatusBadRequest)
		return
	}

	// Chuyển đổi sportId sang số nguyên
	var sportId int
	if _, err := fmt.Sscan(sportIdStr, &sportId); err != nil {
		http.Error(w, "Invalid sportId", http.StatusBadRequest)
		return
	}

	// Xử lý file và sportId
	log.Infof("File: %v", file)
	rq.SportId = sportId

	// Gọi service để tải lên hình ảnh
	err = classService.UploadImageService(imgData, rq.SportId, ctx)
	if err != nil {
		log.WithError(err).Errorf("Error upload image")
		http.Error(w, "Error upload image", http.StatusInternalServerError)
		return
	}

	// Phản hồi thành công
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Image uploaded successfully"))
}

func createSchema(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD POST", "create schema")
	logger.Infof("this API is used to create schema")

	// Kiểm tra phương thức
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Đọc thông tin về schema từ form
	var rq api_request.CreateSchemaRequest
	if err := json.NewDecoder(r.Body).Decode(&rq); err != nil {
		logger.WithError(err).Warningf("Error unmarshaling data from request")
		http.Error(w, "Invalid format", http.StatusBadRequest)
		return
	}

	// Gọi service để tạo schema
	err := classService.CreateSchemaService(ctx, rq)
	if err != nil {
		log.WithError(err).Errorf("Error create schema")
		http.Error(w, "Error create schema", http.StatusInternalServerError)
		return
	}

	// Phản hồi thành công
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Schema created successfully"))
}

func handleCallbackOAuth2(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD GET", "callback oauth2")
	logger.Infof("this API is used to callback oauth2")

	// Kiểm tra phương thức
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Lấy code từ query
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "code is required", http.StatusBadRequest)
		return
	}

	// Phản hồi thành công
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Callback oauth2 successfully"))
}

func handleQueue(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD GET", "callback oauth2")
	logger.Infof("this API is used to callback oauth2")

	// Nâng cấp kết nối HTTP thành WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade:", err)
		return
	}
	defer conn.Close()

	// Giả lập vị trí hàng đợi (bắt đầu từ vị trí 1030)
	position := 5

	// Tạo một vòng lặp để gửi thông báo vị trí mỗi 5 giây
	for position > 0 {
		// Gửi vị trí hiện tại cho client
		err := conn.WriteJSON(map[string]interface{}{
			"position": position,
		})
		if err != nil {
			log.Println("Write error:", err)
			break
		}

		// Giảm vị trí sau mỗi lần gửi (giả lập)
		position--

		// Dừng khi đến vị trí đầu tiên
		if position == 0 {
			err := conn.WriteJSON(map[string]interface{}{
				"message": "Đã đến lượt bạn đặt vé!",
			})
			if err != nil {
				log.Println("Write error:", err)
				break
			}
			break
		}

		// Đợi 5 giây trước khi gửi cập nhật tiếp theo
		time.Sleep(5 * time.Second)
	}
}

func handlePushFileToQueue(w http.ResponseWriter, r *http.Request) {
	ctx := apm.DetachedContext(r.Context())
	logger := GetLoggerWithContext(ctx).WithField("METHOD POST", "push file to queue")
	logger.Infof("this API is used to push file to queue")

	// Kiểm tra phương thức
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Lấy file từ body request
	file, header, err := r.FormFile("file")
	if err != nil {
		logger.WithError(err).Warningf("Error retrieving file from request")
		http.Error(w, "Error retrieving file from request", http.StatusBadRequest)
		return
	}

	logger.Infof("MIME type, File name: %s", header.Filename)

	fileData, err := io.ReadAll(file)
	if err != nil {
		logger.WithError(err).Error("Error reading file data")
		http.Error(w, "Error reading file data", http.StatusInternalServerError)
		return
	}

	mimeType := http.DetectContentType(fileData)
	logger.Infof("MIME type: %s", mimeType)

	defer file.Close()

	// Push file to queue
	data, err := rabbitmq.RabbitMQPublisher(fileData, ctx, mimeType, classService, header.Filename)
	if err != nil {
		log.WithError(err).Errorf("Error push file to queue")
		http.Error(w, "Error push file to queue", http.StatusInternalServerError)
		return
	}

	// Phản hồi thành công
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&data)
}
