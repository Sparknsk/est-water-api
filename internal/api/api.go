package api

import (
	"context"
	"github.com/ozonmp/est-water-api/internal/model"
	"github.com/ozonmp/est-water-api/internal/repo"
	pb "github.com/ozonmp/est-water-api/pkg/est-water-api"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	totalWaterNotFound = promauto.NewCounter(prometheus.CounterOpts{
		Name: "est_water_api_water_not_found_total",
		Help: "Total number of waters that were not found",
	})
)

//go:generate mockgen -destination=../mocks/api_repo_mock.go -package=mocks github.com/ozonmp/est-water-api/internal/repo Repo
type Repo interface {
	DescribeWater(ctx context.Context, waterID uint64) (*model.Water, error)
	CreateWater(ctx context.Context, water *model.Water) error
	ListWaters(ctx context.Context) ([]model.Water, error)
	RemoveWater(ctx context.Context, waterID uint64) error
}

type waterAPI struct {
	pb.UnimplementedEstWaterApiServiceServer
	repo repo.Repo
}

// NewWaterAPI returns api of est-water-api service
func NewWaterAPI(r repo.Repo) pb.EstWaterApiServiceServer {
	return &waterAPI{repo: r}
}

func modelWaterToProtobufWater(water *model.Water) *pb.Water {
	return &pb.Water{
		Id: water.Id,
		Name: water.Name,
		Model: water.Model,
		Manufacturer: water.Manufacturer,
		Material: water.Material,
		Speed: water.Speed,
	}
}
