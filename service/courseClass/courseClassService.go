package courseClass

import (
	"context"
	batman "education-website"
	"education-website/api/request"
	"education-website/api/response/sports"
	"education-website/client"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type classService struct {
	classStore  batman.ClassStore
	flashClient client.FlashClient
}

type ClassServiceCfg struct {
	ClassStore  batman.ClassStore
	FlashClient client.FlashClient
}

func NewClassService(cfg ClassServiceCfg) batman.ClassService {
	return classService{
		classStore:  cfg.ClassStore,
		flashClient: cfg.FlashClient,
	}
}

func (c classService) GetSportsService(ctx context.Context) ([]sports.SportsResponse, error) {
	log.Infof("Start to get sports")
	var responseSport []sports.SportsResponse
	sportListEntities, err := c.classStore.GetSportsStore(ctx)
	if err != nil {
		log.WithError(err).Errorf("Error getting sports")
		return nil, err
	}
	for _, sport := range sportListEntities {
		responseSport = append(responseSport, sports.SportsResponse{
			SportId:   sport.SportId,
			SportName: sport.SportName,
			SportUlr:  sport.SportUlr,
		})
	}

	return responseSport, nil
}

func (c classService) UploadImageService(imageFile []byte, sportId int, ctx context.Context) error {
	return c.classStore.UploadImageStore(imageFile, sportId, ctx)
}

func (c classService) CreateSchemaService(ctx context.Context, request request.CreateSchemaRequest) error {
	city, err := c.classStore.GetCityStore(ctx, request.CityId)
	if city == "" {
		log.Errorf("City not found")
		return fmt.Errorf("City not found")
	}
	if err != nil {
		log.WithError(err).Errorf("Error getting city")
		return err
	}
	createSchemaSql := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %v", request.CityId)
	err = c.classStore.CreateSchemaStore(ctx, createSchemaSql, request)
	if err != nil {
		log.WithError(err).Errorf("Error creating schema")
		return err
	}
	return nil
}
