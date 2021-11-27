package api

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ozonmp/est-water-api/internal/logger"
	"github.com/ozonmp/est-water-api/internal/metric"
	pb "github.com/ozonmp/est-water-api/pkg/est-water-api"
)

func (w *waterAPI) CreateWaterV1 (
	ctx context.Context,
	req *pb.CreateWaterV1Request,
) (*pb.CreateWaterV1Response, error) {

	if err := req.Validate(); err != nil {
		logger.ErrorKV(ctx, "CreateWaterV1() validation error",
			"err", errors.Wrapf(err, "req.Validate() failed with %v", req),
		)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	water, err := w.waterService.CreateWater(ctx, req.Name, req.Model, req.Manufacturer, req.Material, req.Speed)
	if err != nil {
		logger.ErrorKV(ctx, "CreateWaterV1() unable to create",
			"err", errors.Wrapf(err, "waterService.CreateWater() failed with %v", req),
		)

		return nil, status.Error(codes.Internal, "unable to create water entity")
	}

	metric.IncTotalWaterState(metric.StateCreate)

	return &pb.CreateWaterV1Response{
		Water: water.ModelWaterToProtobufWater(),
	}, nil
}
