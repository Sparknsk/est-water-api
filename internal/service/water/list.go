package water_service

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ozonmp/est-water-api/internal/model"
)

func (s *waterService) ListWaters(ctx context.Context, limit uint64, offset uint64) ([]model.Water, error) {
	waters, err := s.waterRepository.List(ctx, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, "waterRepository.List()")
	}

	return waters, err
}
