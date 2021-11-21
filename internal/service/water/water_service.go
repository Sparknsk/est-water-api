package water_service

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/ozonmp/est-water-api/internal/model"
)

var (
	WaterNotFound = errors.New("water entity not found")
)

//go:generate mockgen -destination=../../mocks/service_mock.go -package=mocks github.com/ozonmp/est-water-api/internal/service/water Service
type Service interface {
	DescribeWater(ctx context.Context, waterId uint64) (*model.Water, error)
	CreateWater(ctx context.Context, waterName string, waterModel string, waterMaterial string, waterManufacturer string, waterSpeed uint32) (*model.Water, error)
	ListWaters(ctx context.Context, limit uint64, offset uint64) ([]model.Water, error)
	RemoveWater(ctx context.Context, waterId uint64) error
	UpdateWater(ctx context.Context, waterId uint64, waterName string, waterModel string, waterManufacturer string, waterMaterial string, waterSpeed uint32) (*model.Water, error)
}

type Repo interface {
	Get(ctx context.Context, waterId uint64) (*model.Water, error)
	Create(ctx context.Context, water *model.Water) error
	List(ctx context.Context, limit uint64, offset uint64) ([]model.Water, error)
	Remove(ctx context.Context, waterId uint64) error
	Update(ctx context.Context, water *model.Water) error
}

type EventRepo interface {
	Lock(ctx context.Context, n uint64) ([]model.WaterEvent, error)
	Unlock(ctx context.Context, eventIDs []uint64) error
	Add(ctx context.Context, events []model.WaterEvent) error
	Remove(ctx context.Context, eventIDs []uint64) error
}

type waterService struct {
	db *sqlx.DB
	waterRepository Repo
	waterEventRepository EventRepo
}

func NewService(db *sqlx.DB, waterRepository Repo, waterEventRepository EventRepo) Service {
	return &waterService{
		db: db,
		waterRepository: waterRepository,
		waterEventRepository: waterEventRepository,
	}
}

