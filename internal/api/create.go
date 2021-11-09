package api

import (
	"context"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/ozonmp/est-water-api/pkg/est-water-api"
)

func (w *waterAPI) CreateWaterV1 (
	ctx context.Context,
	req *pb.CreateWaterV1Request,
) (*pb.CreateWaterV1Response, error) {

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("CreateWaterV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	water, err := w.waterService.CreateWater(ctx, req.Name, req.Model, req.Manufacturer, req.Material, req.Speed)
	if err != nil {
		log.Error().Err(err).Msg("CreateWaterV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateWaterV1Response{
		Water: modelWaterToProtobufWater(water),
	}, nil
}
