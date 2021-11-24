package api

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ozonmp/est-water-api/internal/logger"
	pb "github.com/ozonmp/est-water-api/pkg/est-water-api"
)

func (w *waterAPI) ListWatersV1(
	ctx context.Context,
	req *pb.ListWatersV1Request,
) (*pb.ListWatersV1Response, error) {

	if err := req.Validate(); err != nil {
		logger.ErrorKV(ctx, "ListWatersV1() validation error",
			"err", errors.Wrapf(err, "req.Validate() failed with %v", req),
		)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	waters, err := w.waterService.ListWaters(ctx, req.Limit, req.Offset)
	if err != nil {
		logger.ErrorKV(ctx, "ListWatersV1() unable to list",
			"err", errors.Wrapf(err, "waterService.ListWaters() failed with %v", req),
		)

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
