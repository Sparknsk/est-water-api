package api

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ozonmp/est-water-api/internal/service/water"
	pb "github.com/ozonmp/est-water-api/pkg/est-water-api"
)

func (w *waterAPI) DescribeWaterV1(
	ctx context.Context,
	req *pb.DescribeWaterV1Request,
) (*pb.DescribeWaterV1Response, error) {

	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	water, err := w.waterService.DescribeWater(ctx, req.WaterId)
	if err != nil {
		if errors.Is(err, water_service.WaterNotFound) {
			totalWaterNotFound.Inc()

			return nil, status.Error(codes.NotFound, fmt.Sprintf("water entity (id %d) not found", req.WaterId))
		}

		log.Error().Err(errors.Wrap(err, "DescribeWaterV1() failed")).Msg("DescribeWaterV1() unable to get")

		return nil, status.Error(codes.Internal, "unable to get water entity")
	}

	return &pb.DescribeWaterV1Response{
		Water: modelWaterToProtobufWater(water),
	}, nil
}
