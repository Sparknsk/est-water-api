package api

import (
	"context"

	"github.com/pkg/errors"
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
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := w.waterService.RemoveWater(ctx, req.WaterId); err != nil {
		log.Error().Err(errors.Wrap(err, "RemoveWaterV1() failed")).Msg("RemoveWaterV1() unable to remove")

		return nil, status.Error(codes.Internal, "unable to remove water entity")
	}

	return &pb.RemoveWaterV1Response{}, nil
}
