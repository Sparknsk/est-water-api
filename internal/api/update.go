package api

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ozonmp/est-water-api/internal/logger"
	"github.com/ozonmp/est-water-api/internal/service/water"
	pb "github.com/ozonmp/est-water-api/pkg/est-water-api"
)

func (w *waterAPI) UpdateWaterV1 (
	ctx context.Context,
	req *pb.UpdateWaterV1Request,
) (*pb.UpdateWaterV1Response, error) {

	if err := req.Validate(); err != nil {
		logger.ErrorKV(ctx, "UpdateWaterV1() validation error",
			"err", errors.Wrapf(err, "req.Validate() failed with %v", req),
		)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	water, err := w.waterService.UpdateWater(ctx, req.WaterId, req.Name, req.Speed)
	if err != nil {
		if errors.Is(err, water_service.WaterNotFound) {
			totalWaterNotFound.Inc()

			return nil, status.Error(codes.NotFound, fmt.Sprintf("water entity (id %d) not found", req.WaterId))
		}

		logger.ErrorKV(ctx, "UpdateWaterV1() unable to remove",
			"err", errors.Wrapf(err, "waterService.UpdateWater() failed with %v", req),
		)

		return nil, status.Error(codes.Internal, "unable to update water entity")
	}

	return &pb.UpdateWaterV1Response{
		Water: modelWaterToProtobufWater(water),
	}, nil
}
