package api

import (
	"context"

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
		log.Error().Err(err).Msg("ListWatersV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	waters, err := w.repo.ListWaters(ctx)
	if err != nil {
		log.Error().Err(err).Msg("ListWatersV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Msg("ListWatersV1 - success")

	var watersPb []*pb.Water
	for _, water := range waters {
		watersPb = append(watersPb, modelWaterToProtobufWater(&water))
	}

	return &pb.ListWatersV1Response{
		Waters: watersPb,
	}, nil
}
