package client

import (
	"education-website/client/response"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type FlashClient interface {
	GetCourseRevenueByCourseId(courseId string, err error) (*response.GetCourseRevenueByCourseIdResponse, error)
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
