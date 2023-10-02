package courseClass

import (
	"context"
	batman "education-website"
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

func (c classService) GetClassInformationByClassName(request batman.ClassRequest, ctx context.Context) (batman.ClassInfoResponse, error) {
	log.Info("Get class information by class name")

	rs, err := c.classStore.GetClassInformationByClassNameStore(request, ctx)
	if err != nil {
		log.WithError(err).Errorf("Error getting clas information by class Name")
		return nil, err
	}

	return batman.ClassInfoResponse{
		ClassName:  rs.ClassName,
		CourseName: rs.CourseName,
		StudyDate:
	}, nil
}
