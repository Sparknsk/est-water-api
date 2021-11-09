package water_service

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/ozonmp/est-water-api/internal/model"
)

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
