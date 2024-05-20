package client

import (
	"bytes"
	"context"
	"education-website/api/request"
	"education-website/client/response"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"path"
	"time"
)

type FlashClient interface {
	GetCourseRevenueByCourseId(courseId string, err error) (*response.GetCourseRevenueByCourseIdResponse, error)
	GetRevenueByYear(arr []request.GetCourseRevenueByCourseIdRequest, year string, ctx context.Context) (*response.YearReportRevenueResponse, error)
	GetAllAvailableCourseFee() ([]response.CoursesFeeResponse, error)
	GetPaymentStatusByCourseId(courseId request.PaymentStatusRequest, ctx context.Context) ([]response.PaymentStatusByCourseIdResponse, error)
}

type flashClient struct {
	httpClient                 *http.Client
	root                       string
	getCourseRevenueByCourseId string
	getYearlyRevenue           string
}

type FlashClientCfg struct {
	Root                       string `yaml:"root"`
	GetCourseRevenueByCourseId string `yaml:"get-course-revenue-by-course-id"`
	GetYearlyRevenue           string `yaml:"get-yearly-revenue"`
}

func NewFlashClient(cfg FlashClientCfg) FlashClient {
	return flashClient{
		root:                       cfg.Root,
		getCourseRevenueByCourseId: cfg.GetCourseRevenueByCourseId,
		getYearlyRevenue:           cfg.GetYearlyRevenue,
		httpClient: &http.Client{
			Timeout: time.Second * 10, // Example timeout configuration
		},
	}
}

func (f flashClient) GetCourseRevenueByCourseId(courseId string, err error) (*response.GetCourseRevenueByCourseIdResponse, error) {
	log.Infof("Start to get revenue for course %s", courseId)
	url := fmt.Sprintf("%s/%s?param=%s", f.root, f.getCourseRevenueByCourseId, courseId)
	url = "http://34.100.254.97:8083/b/v1/flash/goodbye?param=123"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Perform HTTP request
	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode response body
	var responseBody response.GetCourseRevenueByCourseIdResponse
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return nil, err
	}

	return &responseBody, nil
}

func (f flashClient) GetRevenueByYear(arr []request.GetCourseRevenueByCourseIdRequest, year string, ctx context.Context) (*response.YearReportRevenueResponse, error) {
	log.Infof("Start to get revenue for all course")

	// Chuyển đổi mảng arr thành JSON
	requestBody, err := json.Marshal(arr)
	if err != nil {
		return nil, err
	}

	// Tạo một struct URL từ URL root
	baseURL, err := url.Parse(f.root)
	if err != nil {
		return nil, err
	}

	// Thêm các tham số vào URL
	params := url.Values{}
	params.Set("year", year) // Thêm tham số 'year' vào URL
	baseURL.Path = path.Join(baseURL.Path, f.getYearlyRevenue)
	baseURL.RawQuery = params.Encode()

	// Tạo HTTP request với phần thân là dữ liệu JSON của mảng arr
	req, err := http.NewRequest(http.MethodPost, baseURL.String(), bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	// Thiết lập loại nội dung là JSON
	req.Header.Set("Content-Type", "application/json")

	// Thực hiện HTTP request
	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Kiểm tra mã trạng thái của response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode response body
	var responseBody response.YearReportRevenueResponse
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return nil, err
	}

	return &responseBody, nil
}

func (f flashClient) GetAllAvailableCourseFee() ([]response.CoursesFeeResponse, error) {
	log.Infof("Start to get revenue for all course")
	url := fmt.Sprintf("%s/%s?param=%s", f.root, f.getCourseRevenueByCourseId)
	url = "http://34.100.254.97:8083/b/v1/flash/goodbye?param=123"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Perform HTTP request
	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode response body
	var responseBody []response.CoursesFeeResponse
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return nil, err
	}

	return responseBody, nil
}

func (f flashClient) GetPaymentStatusByCourseId(courseId request.PaymentStatusRequest, ctx context.Context) ([]response.PaymentStatusByCourseIdResponse, error) {
	log.Infof("Start to get payment status for course %s", courseId)
	url := fmt.Sprintf("%s/%s?param=%s", f.root, f.getCourseRevenueByCourseId)
	url = "http://34.100.254.97:8083/b/v1/flash/paymentStatus"

	requestBody, err := json.Marshal(courseId)
	if err != nil {
		log.WithError(err).Errorf("Error when marshalling request body for course ID: %s", courseId)
		return nil, err
	}

	// Tạo yêu cầu HTTP với body
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.WithError(err).Errorf("Error when creating new request for url: %s", url)
		return nil, err
	}

	// Đặt header cho yêu cầu
	req.Header.Set("Content-Type", "application/json")

	// Perform HTTP request
	resp, err := f.httpClient.Do(req)
	if err != nil {
		log.WithError(err).Errorf("Error when performing HTTP request")
		return nil, err
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode response body
	var responseBody []response.PaymentStatusByCourseIdResponse
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		log.WithError(err).Errorf("Error when decoding response body for url: %s", url)
		return nil, err
	}
	return responseBody, nil
}
