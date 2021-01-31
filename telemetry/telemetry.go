package telemetry

import (
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const InstrumentationName = "github.com/eran-levy/tokenizer-gophercon"

var TracerConfigIsNotValid = errors.New("provided tracer configuration is invalid")

type TracerConfig struct {
	ApplicationID string
	ServiceName   string
	AgentEndpoint string
}

type Tracer struct {
	Tracer trace.Tracer
	Config TracerConfig
}

func New(config TracerConfig) (Tracer, func(), error) {
	const appId = "app_id"
	if !isConfigValid(config) {
		return Tracer{}, nil, TracerConfigIsNotValid
	}
	t := otel.Tracer(InstrumentationName)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	flush, err := jaeger.InstallNewPipeline(
		jaeger.WithAgentEndpoint(config.AgentEndpoint),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: config.ServiceName,
			Tags: []label.KeyValue{
				label.String(appId, config.ApplicationID),
			},
		}),
		jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
	)
	return Tracer{Tracer: t, Config: config}, flush, err
}

func isConfigValid(config TracerConfig) bool {
	if len(config.ApplicationID) == 0 {
		return false
	}
	if len(config.ServiceName) == 0 {
		return false
	}
	if len(config.AgentEndpoint) == 0 {
		return false
	}
	return true
}
