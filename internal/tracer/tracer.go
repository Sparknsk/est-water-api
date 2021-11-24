package tracer

import (
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"

	"github.com/ozonmp/est-water-api/internal/config"
)

func NewTracer(cfg *config.Config) (io.Closer, error) {
	jaegerAddr := fmt.Sprintf("%s:%d", cfg.Jaeger.Host, cfg.Jaeger.Port)
	cfgTracer := &jaegercfg.Configuration{
		ServiceName: cfg.Jaeger.Service,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: jaegerAddr,
		},
	}
	tracer, closer, err := cfgTracer.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		return nil, errors.Wrap(err, "cfgTracer.NewTracer() failed")
	}
	opentracing.SetGlobalTracer(tracer)

	return closer, nil
}
