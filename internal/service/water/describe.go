package water_service

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ozonmp/est-water-api/internal/model"
)

func (s *waterService) DescribeWater(ctx context.Context, WaterId uint64) (*model.Water, error) {
	water, err := s.waterRepository.Get(ctx, WaterId)
	if err != nil {
		return nil, errors.Wrap(err, "waterRepository.Get() failed")
	}

	if water == nil {
		return nil, WaterNotFound
	}

	return water, nil
}
