package water_service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/ozonmp/est-water-api/internal/model"
)

func (s *waterService) ListWaters(ctx context.Context, limit uint64, offset uint64) ([]model.Water, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "waterService.ListWaters()")
	defer span.Finish()
	span.LogKV(
		"event", "service list water",
		"limit", limit,
		"offset", offset,
	)

	waters, err := s.waterRepository.List(ctx, limit, offset)
	if err != nil {
		return nil, errors.Wrapf(err, "waterRepository.List() failed with limit=%d, offset=%d", limit, offset)
	}

	return waters, nil
}
