package api

import (
	"context"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/ozonmp/est-water-api/pkg/est-water-api"
)

func (w *waterAPI) UpdateWaterV1 (
	ctx context.Context,
	req *pb.UpdateWaterV1Request,
) (*pb.UpdateWaterV1Response, error) {

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("UpdateWaterV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	water, err := w.waterService.UpdateWater(ctx, req.WaterId, req.Name, req.Speed)
	if err != nil {
		log.Error().Err(err).Msg("UpdateWaterV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	if water == nil {
		totalWaterNotFound.Inc()

		return nil, status.Error(codes.NotFound, "water not found")
	}

	return &pb.UpdateWaterV1Response{
		Water: modelWaterToProtobufWater(water),
	}, nil
}
