package server

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/ozonmp/est-water-api/internal/logger"
)

func (s *GrpcServer) requestInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,) (interface{}, error) {

	logEnabled := false
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		logEnabledStr := md.Get(s.cfg.Logging.HeaderNameForResponseLog)
		if len(logEnabledStr) > 0 {
			logEnabled = true
		}
	}

	var res interface{}
	var err error
	if logEnabled {
		reqDebug := fmt.Sprintf("Request: Method - %s, Data - %v", info.FullMethod, req)

		res, err = handler(ctx, req)

		if err != nil {
			logger.InfoKV(ctx, fmt.Sprintf("%v | Response with error: %v", reqDebug, err))
		} else {
			logger.InfoKV(ctx, fmt.Sprintf("%v | Response with success: %v", reqDebug, res))
		}
	} else {
		res, err = handler(ctx, req)
	}

	return res, err
}