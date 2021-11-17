package water_service

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/ozonmp/est-water-api/internal/model"
)

func (s *waterService) UpdateWater(ctx context.Context, waterId uint64, waterName string, waterSpeed uint32) (*model.Water, error) {
	water, err := s.waterRepository.Get(ctx, waterId)
	if err != nil {
		return nil, errors.Wrap(err, "waterRepository.Get() failed")
	}

	if water == nil {
		return nil, WaterNotFound
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
			return nil, errors.Wrap(err, "db.BeginTxx() failed")
		}

		if err := s.waterRepository.Update(ctx, water); err != nil {
			return nil, errors.Wrap(err, "waterRepository.Update() failed")
		}

		if err := s.waterEventRepository.Add(ctx, waterEvents); err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, errors.Wrap(err, "tx.Rollback() failed")
			}
			return nil, errors.Wrap(err, "waterEventRepository.Add() failed")
		}

		if err := tx.Commit(); err != nil {
			return nil, errors.Wrap(err, "tx.Commit() failed")
		}
	}

	return water, nil
}
