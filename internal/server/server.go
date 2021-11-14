package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"github.com/ozonmp/est-water-api/internal/api"
	er "github.com/ozonmp/est-water-api/internal/app/repo"
	"github.com/ozonmp/est-water-api/internal/config"
	"github.com/ozonmp/est-water-api/internal/logger"
	"github.com/ozonmp/est-water-api/internal/repo"
	"github.com/ozonmp/est-water-api/internal/service/water"
	pb "github.com/ozonmp/est-water-api/pkg/est-water-api"
)

// GrpcServer is gRPC server
type GrpcServer struct {
	cfg *config.Config
	db *sqlx.DB
	batchSize uint
}

// NewGrpcServer returns gRPC server with supporting of batch listing
func NewGrpcServer(cfg *config.Config, db *sqlx.DB, batchSize uint) *GrpcServer {
	return &GrpcServer{
		cfg: cfg,
		db: db,
		batchSize: batchSize,
	}
}

// Start method runs server
func (s *GrpcServer) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gatewayAddr := fmt.Sprintf("%s:%v", s.cfg.Rest.Host, s.cfg.Rest.Port)
	grpcAddr := fmt.Sprintf("%s:%v", s.cfg.Grpc.Host, s.cfg.Grpc.Port)
	metricsAddr := fmt.Sprintf("%s:%v", s.cfg.Metrics.Host, s.cfg.Metrics.Port)

	gatewayServer := createGatewayServer(grpcAddr, gatewayAddr)

	go func() {
		logger.InfoKV(ctx, fmt.Sprintf("Gateway server is running on %s", gatewayAddr))
		if err := gatewayServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.ErrorKV(ctx, "Failed running gateway server",
				"err", errors.Wrap(err, "gatewayServer.ListenAndServe() failed"),
			)
			cancel()
		}
	}()

	metricsServer := createMetricsServer(s.cfg)

	go func() {
		logger.InfoKV(ctx, fmt.Sprintf("Metrics server is running on %s", metricsAddr))
		if err := metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.ErrorKV(ctx, "Failed running metrics server",
				"err", errors.Wrap(err, "metricsServer.ListenAndServe() failed"),
			)
			cancel()
		}
	}()

	isReady := &atomic.Value{}
	isReady.Store(false)

	statusServer := createStatusServer(s.cfg, isReady)

	go func() {
		statusAdrr := fmt.Sprintf("%s:%v", s.cfg.Status.Host, s.cfg.Status.Port)
		logger.InfoKV(ctx, fmt.Sprintf("Status server is running on %s", statusAdrr))
		if err := statusServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.ErrorKV(ctx, "Failed running starus server",
				"err", errors.Wrap(err, "statusServer.ListenAndServe() failed"),
			)
		}
	}()

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer l.Close()

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: time.Duration(s.cfg.Grpc.MaxConnectionIdle) * time.Minute,
			Timeout:           time.Duration(s.cfg.Grpc.Timeout) * time.Second,
			MaxConnectionAge:  time.Duration(s.cfg.Grpc.MaxConnectionAge) * time.Minute,
			Time:              time.Duration(s.cfg.Grpc.Timeout) * time.Minute,
		}),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			s.requestInterceptor,
			s.loggerLevelInterceptor,
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_opentracing.UnaryServerInterceptor(),
			grpcrecovery.UnaryServerInterceptor(),
		)),
	)

	repository := repo.NewRepo(s.db, s.batchSize)
	eventRepository := er.NewEventRepo(s.db)
	service := water_service.NewService(s.db, repository, eventRepository)

	pb.RegisterEstWaterApiServiceServer(grpcServer, api.NewWaterAPI(service))
	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(grpcServer)

	go func() {
		logger.InfoKV(ctx, fmt.Sprintf("GRPC Server is listening on: %s", grpcAddr))
		if err := grpcServer.Serve(l); err != nil {
			logger.ErrorKV(ctx, "Failed running gRPC server",
				"err", errors.Wrap(err, "grpcServer.Serve() failed"),
			)
		}
	}()

	go func() {
		time.Sleep(2 * time.Second)
		isReady.Store(true)
		logger.InfoKV(ctx, "The service is ready to accept requests")
	}()

	if s.cfg.Logging.IsDebug() {
		reflection.Register(grpcServer)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		logger.InfoKV(ctx, fmt.Sprintf("signal.Notify: %v", v))
	case done := <-ctx.Done():
		logger.InfoKV(ctx, fmt.Sprintf("ctx.Done: %v", done))
	}

	isReady.Store(false)

	if err := gatewayServer.Shutdown(ctx); err != nil {
		logger.ErrorKV(ctx, "gatewayServer.Shutdown",
			"err", errors.Wrap(err, "gatewayServer.Shutdown() failed"),
		)
	} else {
		logger.InfoKV(ctx, "gatewayServer shut down correctly")
	}

	if err := statusServer.Shutdown(ctx); err != nil {
		logger.ErrorKV(ctx, "statusServer.Shutdown",
			"err", errors.Wrap(err, "statusServer.Shutdown() failed"),
		)
	} else {
		logger.InfoKV(ctx, "statusServer shut down correctly")
	}

	if err := metricsServer.Shutdown(ctx); err != nil {
		logger.ErrorKV(ctx, "metricsServer.Shutdown",
			"err", errors.Wrap(err, "metricsServer.Shutdown() failed"),
		)
	} else {
		logger.InfoKV(ctx, "metricsServer shut down correctly")
	}

	grpcServer.GracefulStop()
	logger.InfoKV(ctx, "grpcServer shut down correctly")

	return nil
}

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
