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

func (w *waterAPI) DescribeWaterV1(
	ctx context.Context,
	req *pb.DescribeWaterV1Request,
) (*pb.DescribeWaterV1Response, error) {

	if err := req.Validate(); err != nil {
		logger.ErrorKV(ctx, "DescribeWaterV1() validation error",
			"err", errors.Wrapf(err, "req.Validate() failed with %v", req),
		)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	water, err := w.waterService.DescribeWater(ctx, req.WaterId)
	if err != nil {
		if errors.Is(err, water_service.WaterNotFound) {
			metric.IncTotalWaterNotFound()

			return nil, status.Error(codes.NotFound, fmt.Sprintf("water entity (id %d) not found", req.WaterId))
		}

		logger.ErrorKV(ctx, "DescribeWaterV1() unable to get",
			"err", errors.Wrapf(err, "waterService.DescribeWater() failed with %v", req),
		)

		return nil, status.Error(codes.Internal, "unable to get water entity")
	}

	return &pb.DescribeWaterV1Response{
		Water: modelWaterToProtobufWater(water),
	}, nil
}
