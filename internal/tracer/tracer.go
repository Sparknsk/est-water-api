package tracer

import (
	"context"
	"github.com/pkg/errors"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"

	"github.com/ozonmp/est-water-api/internal/config"
	"github.com/ozonmp/est-water-api/internal/logger"
)

func NewTracer(cfg *config.Config) (io.Closer, error) {
	ctx := context.Background()

	cfgTracer := &jaegercfg.Configuration{
		ServiceName: cfg.Jaeger.Service,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: cfg.Jaeger.Host + cfg.Jaeger.Port,
		},
	}
	tracer, closer, err := cfgTracer.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		logger.ErrorKV(ctx, "failed init jaeger",
			"err", errors.Wrap(err, "cfgTracer.NewTracer() failed"),
		)

		return nil, err
	}
	opentracing.SetGlobalTracer(tracer)
	logger.InfoKV(ctx, "Traces started")

	return closer, nil
}
