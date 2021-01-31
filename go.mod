module github.com/eran-levy/tokenizer-gophercon

go 1.15

require (
	github.com/gin-contrib/pprof v1.3.0
	github.com/gin-gonic/gin v1.6.3
	github.com/google/uuid v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/viper v1.7.1
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.16.0
	go.opentelemetry.io/contrib/instrumentation/host v0.16.0
	go.opentelemetry.io/contrib/instrumentation/runtime v0.16.0
	go.opentelemetry.io/otel v0.16.0
	go.opentelemetry.io/otel/exporters/metric/prometheus v0.16.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.16.0
	go.opentelemetry.io/otel/sdk v0.16.0
	go.uber.org/zap v1.13.0
)
