package api

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/ozonmp/est-water-api/pkg/est-water-api"
)

func (w *waterAPI) RemoveWaterV1(
	ctx context.Context,
	req *pb.RemoveWaterV1Request,
) (*pb.RemoveWaterV1Response, error) {

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("RemoveWaterV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := w.repo.RemoveWater(ctx, req.WaterId); err != nil {
		log.Error().Err(err).Msg("RemoveWaterV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RemoveWaterV1Response{}, nil
}
