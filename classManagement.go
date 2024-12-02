package education_website

import (
	"context"
	"education-website/api/request"
	"education-website/api/response/sports"
	sports2 "education-website/entity/sports"
)

type ClassService interface {
	GetSportsService(ctx context.Context) ([]sports.SportsResponse, error)
	UploadImageService(imageFile []byte, sportId int, ctx context.Context) error
	CreateSchemaService(ctx context.Context, request request.CreateSchemaRequest) error
}

type ClassStore interface {
	GetSportsStore(ctx context.Context) ([]sports2.SportsEntity, error)
	UploadImageStore(imageFile []byte, sportId int, ctx context.Context) error
	CreateSchemaStore(ctx context.Context, createSchemaSql string, request request.CreateSchemaRequest) error
	GetCityStore(ctx context.Context, cityId int) (string, error)
}
