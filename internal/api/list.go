package api

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/ozonmp/est-water-api/pkg/est-water-api"
)

func (w *waterAPI) ListWatersV1(
	ctx context.Context,
	req *pb.ListWatersV1Request,
) (*pb.ListWatersV1Response, error) {

	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	waters, err := w.waterService.ListWaters(ctx, req.Limit, req.Offset)
	if err != nil {
		log.Error().Err(errors.Wrap(err, "ListWatersV1() failed")).Msg("ListWatersV1() unable to list")

		return nil, status.Error(codes.Internal, "unable to list water entity")
	}

	var watersPb []*pb.Water
	for _, water := range waters {
		watersPb = append(watersPb, modelWaterToProtobufWater(&water))
	}

	return &pb.ListWatersV1Response{
		Waters: watersPb,
	}, nil
}
