package server

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/ozonmp/est-water-api/internal/logger"
)

func (s *GrpcServer) loggerLevelInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,) (interface{}, error) {

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		levelStr := md.Get(s.cfg.Logging.HeaderNameForRequestLevel)
		if len(levelStr) > 0 {
			if newLogLevel, ok := logger.LevelFromString(levelStr[0]); ok {
				logger.InfoKV(ctx, fmt.Sprintf("Set %s log level for request", levelStr[0]))

				newLogger := logger.CloneWithLevel(ctx, newLogLevel)
				ctx = logger.AttachLogger(ctx, newLogger)
			}
		}
	}

	return handler(ctx, req)
}