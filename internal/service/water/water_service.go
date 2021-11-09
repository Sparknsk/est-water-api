package water_service

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ozonmp/est-water-api/internal/model"
	"github.com/pkg/errors"
)

//go:generate mockgen -destination=../../mocks/service_mock.go -package=mocks github.com/ozonmp/est-water-api/internal/service/water Service
type Service interface {
	DescribeWater(ctx context.Context, waterId uint64) (*model.Water, error)
	CreateWater(ctx context.Context, waterName string, waterModel string, waterMaterial string, waterManufacturer string, waterSpeed uint32) (*model.Water, error)
	ListWaters(ctx context.Context, limit uint64, offset uint64) ([]model.Water, error)
	RemoveWater(ctx context.Context, waterId uint64) error
	UpdateWater(ctx context.Context, waterId uint64, waterName string, waterSpeed uint32) (*model.Water, error)
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

func (s *waterService) DescribeWater(ctx context.Context, WaterId uint64) (*model.Water, error) {
	water, err := s.waterRepository.Get(ctx, WaterId)
	if err != nil {
		return nil, errors.Wrap(err, "waterRepository.Get()")
	}

	return water, err
}

func (s *waterService) CreateWater(ctx context.Context, waterName string, waterModel string, waterMaterial string, waterManufacturer string, waterSpeed uint32) (*model.Water, error) {

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "db.BeginTxx()")
	}

	ts := time.Now().UTC()
	water := model.Water{
		Name: waterName,
		Model: waterModel,
		Material: waterMaterial,
		Manufacturer: waterManufacturer,
		Speed: waterSpeed,
		CreatedAt: &ts,
	}
	if err := s.waterRepository.Create(ctx, &water); err != nil {
		return nil, errors.Wrap(err, "waterRepository.Create()")
	}

	waterEvent := model.WaterEvent{
		WaterId: water.Id,
		Type: model.Created,
		Status: model.Unlocked,
		Entity: &water,
		CreatedAt: &ts,
	}
	if err := s.waterEventRepository.Add(ctx, []model.WaterEvent{waterEvent}); err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, errors.Wrap(err, "tx.Rollback()")
		}
		return nil, errors.New("waterEventRepository.Add()")
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "tx.Commit()")
	}

	return &water, nil
}

func (s *waterService) ListWaters(ctx context.Context, limit uint64, offset uint64) ([]model.Water, error) {
	waters, err := s.waterRepository.List(ctx, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, "waterRepository.List()")
	}

	return waters, err
}

func (s *waterService) RemoveWater(ctx context.Context, waterId uint64) error {
	water, err := s.waterRepository.Get(ctx, waterId)
	if err != nil {
		return errors.Wrap(err, "waterRepository.Get()")
	}

	if water == nil {
		return nil
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "db.BeginTxx()")
	}

	if err := s.waterRepository.Remove(ctx, waterId); err != nil {
		return errors.Wrap(err, "waterRepository.Remove()")
	}

	ts := time.Now().UTC()
	waterEvent := model.WaterEvent{
		WaterId: water.Id,
		Type: model.Removed,
		Status: model.Unlocked,
		Entity: water,
		CreatedAt: &ts,
	}
	if err := s.waterEventRepository.Add(ctx, []model.WaterEvent{waterEvent}); err != nil {
		if err := tx.Rollback(); err != nil {
			return errors.Wrap(err, "tx.Rollback()")
		}
		return errors.New("waterEventRepository.Add()")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "tx.Commit()")
	}

	return nil
}

func (s *waterService) UpdateWater(ctx context.Context, waterId uint64, waterName string, waterSpeed uint32) (*model.Water, error) {
	water, err := s.waterRepository.Get(ctx, waterId)
	if err != nil {
		return nil, errors.Wrap(err, "waterRepository.Get()")
	}

	fmt.Println(water)

	if water == nil {
		return nil, errors.Wrap(err, "Entity not found")
	}

	ts := time.Now().UTC()
	water.UpdatedAt = &ts

	var waterEvents []model.WaterEvent

	if waterName != water.Name {
		water.Name = waterName

		waterCopy := *water
		waterEvents = append(
			waterEvents,
			model.WaterEvent{
				WaterId: water.Id,
				Type: model.UpdatedName,
				Status: model.Unlocked,
				Entity: &waterCopy,
				CreatedAt: &ts,
			},
		)
	}

	if waterSpeed != water.Speed {
		water.Speed = waterSpeed

		waterCopy := *water
		waterEvents = append(
			waterEvents,
			model.WaterEvent{
				WaterId: water.Id,
				Type: model.UpdatedSpeed,
				Status: model.Unlocked,
				Entity: &waterCopy,
				CreatedAt: &ts,
			},
		)
	}

	if len(waterEvents) > 0 {
		tx, err := s.db.BeginTxx(ctx, nil)
		if err != nil {
			return nil, errors.Wrap(err, "db.BeginTxx()")
		}

		if err := s.waterRepository.Update(ctx, water); err != nil {
			return nil, errors.Wrap(err, "waterRepository.Update()")
		}

		if err := s.waterEventRepository.Add(ctx, waterEvents); err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, errors.Wrap(err, "tx.Rollback()")
			}
			return nil, errors.New("waterEventRepository.Add()")
		}

		if err := tx.Commit(); err != nil {
			return nil, errors.Wrap(err, "tx.Commit()")
		}
	}

	return water, nil
}

