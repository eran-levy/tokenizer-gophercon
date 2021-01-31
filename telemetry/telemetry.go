package telemetry

import (
	"context"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const (
	InstrumentationName = "github.com/eran-levy/tokenizer-gophercon"
	globalTxTagKey      = "request.global_tx_id"
	requestIdTagKey     = "request.id"
)

var (
	GlobalTxIdKey          = label.Key(globalTxTagKey)
	ReuqestIdTagKey        = label.Key(requestIdTagKey)
	TracerConfigIsNotValid = errors.New("provided tracer configuration is invalid")
)

type Config struct {
	ApplicationID string
	ServiceName   string
	AgentEndpoint string
}

type Telemetry struct {
	Tracer trace.Tracer
	Config Config
}

var (
	apiRequestCounter metric.Int64Counter
	serviceNameKV     label.KeyValue
	statusKey         = label.Key("status")
)

const (
	SuccessStatusValue = "SUCCESS"
	FailStatusValue    = "FAIL"
)

func New(config Config) (Telemetry, func(), error) {
	const appId = "app_id"
	if !isConfigValid(config) {
		return Telemetry{}, nil, TracerConfigIsNotValid
	}
	meter := otel.Meter(InstrumentationName)
	err := runtime.Start()
	if err != nil {
		return Telemetry{}, nil, errors.Wrap(err, "could not init runtime metrics")
	}
	err = host.Start()
	if err != nil {
		return Telemetry{}, nil, errors.Wrap(err, "could not init host metrics")
	}
	serviceNameKV = label.String("service", config.ServiceName)
	apiRequestCounter = metric.Must(meter).NewInt64Counter("service_http_request_counter")

	t := otel.Tracer(InstrumentationName)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	flush, err := jaeger.InstallNewPipeline(
		//in real life app, you may use the agent WithAgentEndpoint()
		jaeger.WithCollectorEndpoint(config.AgentEndpoint),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: config.ServiceName,
			Tags: []label.KeyValue{
				label.String(appId, config.ApplicationID),
			},
		}),
		//in real life app, you may use sampler instead of the always sampler
		jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
	)
	return Telemetry{Tracer: t, Config: config}, flush, err
}
func GetMeterHandlerToServe() (*prometheus.Exporter, error) {
	exporter, err := prometheus.InstallNewPipeline(prometheus.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "could not install prometheus exporter pipeline")
	}

	return exporter, nil
}
func isConfigValid(config Config) bool {
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

func IncAPIRequestCounter(ctx context.Context, v int64, status string) {
	apiRequestCounter.Add(ctx, v, serviceNameKV, statusKey.String(status))
}
