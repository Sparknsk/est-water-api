package server

import (
	"context"
	"github.com/pkg/errors"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"

	"github.com/ozonmp/est-water-api/internal/logger"
	pb "github.com/ozonmp/est-water-api/pkg/est-water-api"
)

var (
	httpTotalRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_microservice_requests_total",
		Help: "The total number of incoming HTTP requests",
	})
)

func createGatewayServer(grpcAddr, gatewayAddr string) *http.Server {
	ctx := context.Background()

	// Create a client connection to the gRPC Server we just started.
	// This is where the gRPC-Gateway proxies the requests.
	conn, err := grpc.DialContext(
		context.Background(),
		grpcAddr,
		grpc.WithInsecure(),
	)
	if err != nil {
		logger.FatalKV(ctx, "Failed to dial server",
			"err", errors.Wrap(err, "grpc.DialContext() failed"),
		)
	}

	mux := runtime.NewServeMux()
	if err := pb.RegisterEstWaterApiServiceHandler(context.Background(), mux, conn); err != nil {
		logger.FatalKV(ctx, "Failed registration handler",
			"err", errors.Wrap(err, "pb.RegisterEstWaterApiServiceHandler() failed"),
		)
	}

	gatewayServer := &http.Server{
		Addr:    gatewayAddr,
		Handler: tracingWrapper(mux),
	}

	return gatewayServer
}

var grpcGatewayTag = opentracing.Tag{Key: string(ext.Component), Value: "grpc-gateway"}

func tracingWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpTotalRequests.Inc()
		parentSpanContext, err := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(r.Header))
		if err == nil || errors.Is(err, opentracing.ErrSpanContextNotFound) {
			serverSpan := opentracing.GlobalTracer().StartSpan(
				"ServeHTTP",
				ext.RPCServerOption(parentSpanContext),
				grpcGatewayTag,
			)
			r = r.WithContext(opentracing.ContextWithSpan(r.Context(), serverSpan))
			defer serverSpan.Finish()
		}
		h.ServeHTTP(w, r)
	})
}
