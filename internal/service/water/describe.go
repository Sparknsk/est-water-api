package water_service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/ozonmp/est-water-api/internal/model"
)

func (s *waterService) DescribeWater(ctx context.Context, waterId uint64) (*model.Water, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "waterService.DescribeWater()")
	defer span.Finish()
	span.LogKV(
		"event", "service describe water",
		"waterId", waterId,
	)

	water, err := s.waterRepository.Get(ctx, waterId)
	if err != nil {
		return nil, errors.Wrapf(err, "waterRepository.Get() failed with id=%d", waterId)
	}

	if water == nil {
		return nil, WaterNotFound
	}

	return water, nil
}
