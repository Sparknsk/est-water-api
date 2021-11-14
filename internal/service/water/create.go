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
		return nil, errors.Wrap(err, "db.BeginTxx() failed")
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
		return nil, errors.Wrapf(err, "waterRepository.Create() failed with %v", water)
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
			return nil, errors.Wrap(err, "tx.Rollback() failed")
		}
		return nil, errors.Wrapf(err,"waterEventRepository.Add() failed with %v", waterEvent)
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "tx.Commit() failed")
	}

	return &water, nil
}