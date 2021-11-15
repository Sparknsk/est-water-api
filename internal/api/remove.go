package api

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ozonmp/est-water-api/internal/logger"
	"github.com/ozonmp/est-water-api/internal/metric"
	"github.com/ozonmp/est-water-api/internal/service/water"
	pb "github.com/ozonmp/est-water-api/pkg/est-water-api"
)

func (w *waterAPI) RemoveWaterV1(
	ctx context.Context,
	req *pb.RemoveWaterV1Request,
) (*pb.RemoveWaterV1Response, error) {

	if err := req.Validate(); err != nil {
		logger.ErrorKV(ctx, "RemoveWaterV1() validation error",
			"err", errors.Wrapf(err, "req.Validate() failed with %v", req),
		)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := w.waterService.RemoveWater(ctx, req.WaterId); err != nil {
		if errors.Is(err, water_service.WaterNotFound) {
			metric.IncTotalWaterNotFound()

			return nil, status.Error(codes.NotFound, fmt.Sprintf("water entity (id %d) not found", req.WaterId))
		}

		logger.ErrorKV(ctx, "RemoveWaterV1() unable to remove",
			"err", errors.Wrapf(err, "waterService.RemoveWater() failed with %v", req),
		)

		return nil, status.Error(codes.Internal, "unable to remove water entity")
	}

	metric.IncTotalWaterState(metric.StateRemoved)

	return &pb.RemoveWaterV1Response{}, nil
}
