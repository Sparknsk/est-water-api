package water_service

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/ozonmp/est-water-api/internal/model"
)

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
