package water_service

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/ozonmp/est-water-api/internal/model"
)

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