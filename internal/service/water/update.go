package water_service

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/ozonmp/est-water-api/internal/model"
)

func (s *waterService) UpdateWater(ctx context.Context, waterId uint64, waterName string, waterSpeed uint32) (*model.Water, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "waterService.CreateWater()")
	defer span.Finish()
	span.LogKV(
		"event", "service update water",
		"waterId", waterId,
		"waterName", waterName,
		"waterSpeed", waterSpeed,
	)

	water, err := s.waterRepository.Get(ctx, waterId)
	if err != nil {
		return nil, errors.Wrapf(err, "waterRepository.Get() failed with id=%d", waterId)
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
			return nil, errors.Wrapf(err, "waterRepository.Update() failed with %v", water)
		}

		if err := s.waterEventRepository.Add(ctx, waterEvents); err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, errors.Wrap(err, "tx.Rollback() failed")
			}
			return nil, errors.Wrapf(err,"waterEventRepository.Add() failed with %v", waterEvents)
		}

		if err := tx.Commit(); err != nil {
			return nil, errors.Wrap(err, "tx.Commit() failed")
		}
	}

	return water, nil
}
