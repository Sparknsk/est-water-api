package api

import (
	"context"

	"github.com/ozonmp/est-water-api/internal/model"
	pb "github.com/ozonmp/est-water-api/pkg/est-water-api"
)

type Service interface {
	DescribeWater(ctx context.Context, waterId uint64) (*model.Water, error)
	CreateWater(ctx context.Context, waterName string, waterModel string, waterMaterial string, waterManufacturer string, waterSpeed uint32) (*model.Water, error)
	ListWaters(ctx context.Context, limit uint64, offset uint64) ([]model.Water, error)
	RemoveWater(ctx context.Context, waterId uint64) error
	UpdateWater(ctx context.Context, waterId uint64, waterName string, waterModel string, waterManufacturer string, waterMaterial string, waterSpeed uint32) (*model.Water, error)
}

type waterAPI struct {
	pb.UnimplementedEstWaterApiServiceServer
	waterService Service
}

func NewWaterAPI(waterService Service) pb.EstWaterApiServiceServer {
	return &waterAPI{waterService: waterService}
}
