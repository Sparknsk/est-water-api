package repo

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/ozonmp/est-water-api/internal/model"
)

type Repo interface {
	DescribeWater(ctx context.Context, waterID uint64) (*model.Water, error)
	CreateWater(ctx context.Context, water *model.Water) error
	ListWaters(ctx context.Context) ([]model.Water, error)
	RemoveWater(ctx context.Context, waterID uint64) error
}

type repo struct {
	db        *sqlx.DB
	batchSize uint
}

// NewRepo returns Repo interface
func NewRepo(db *sqlx.DB, batchSize uint) Repo {
	return &repo{db: db, batchSize: batchSize}
}

func (r *repo) DescribeWater(ctx context.Context, waterID uint64) (*model.Water, error) {
	return nil, nil
}

func (r *repo) CreateWater(ctx context.Context, water *model.Water) error {
	return nil
}

func (r *repo) ListWaters(ctx context.Context) ([]model.Water, error) {
	return []model.Water{}, nil
}

func (r *repo) RemoveWater(ctx context.Context, waterID uint64) error {
	return nil
}
