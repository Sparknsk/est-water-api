package water_service

import (
	"context"
	"github.com/ozonmp/est-water-api/internal/model"
	"github.com/pkg/errors"
)

func (s *waterService) DescribeWater(ctx context.Context, WaterId uint64) (*model.Water, error) {
	water, err := s.waterRepository.Get(ctx, WaterId)
	if err != nil {
		return nil, errors.Wrap(err, "waterRepository.Get()")
	}

	return water, err
}
