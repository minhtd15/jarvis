package courseClass

import (
	"context"
	batman "education-website"
	api_request "education-website/api/request"
	api_response "education-website/api/response"
	log "github.com/sirupsen/logrus"
	"time"
)

type classService struct {
	classStore batman.ClassStore
}

type ClassServiceCfg struct {
	ClassStore batman.ClassStore
}

func NewClassService(cfg ClassServiceCfg) batman.ClassService {
	return classService{
		classStore: cfg.ClassStore,
	}
}

//func (c classService) AddNewClass(request api_request.NewClassRequest, ctx context.Context) error {
//
//}

func (c classService) GetCourseInformationByClassName(request api_request.CourseInfoRequest, ctx context.Context) (*api_response.CourseInfoResponse, error) {
	log.Info("Get class information by class name")

	// get class information
	rs, err := c.classStore.GetCourseInformationStore(request, ctx)
	if err != nil {
		log.WithError(err).Errorf("Error getting clas information by class Name")
		return nil, err
	}

	var startDate, endDate time.Time
	if rs.StartDate.Valid {
		startDate, err = time.Parse("2006-01-02", rs.StartDate.String)
		if err != nil {
			log.WithError(err).Errorf("Cannot parse start time")
			return nil, err
		}
	}

	if rs.EndDate.Valid {
		endDate, err = time.Parse("2006-01-02", rs.EndDate.String)
		if err != nil {
			log.WithError(err).Errorf("Cannot parse start time")
			return nil, err
		}
	}

	var courseStatus string
	currentDate := time.Now()
	if currentDate.Before(startDate) || currentDate.Equal(startDate) {
		courseStatus = "INACTIVE"
	} else if currentDate.After(endDate) {
		courseStatus = "FINISHED"
	} else {
		courseStatus = "ACTIVE"
	}

	return &api_response.CourseInfoResponse{
		CourseId:      rs.CourseId,
		CourseName:    rs.CourseName,
		MainTeacher:   rs.MainTeacher,
		Room:          rs.Room,
		StartDate:     startDate,
		EndDate:       endDate,
		StudyDays:     rs.StudyDays,
		CourseStatus:  courseStatus,
		TotalSessions: rs.TotalSessions,
	}, nil
}
