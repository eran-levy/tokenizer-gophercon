package http

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
)

func MetricHandler(exp *prometheus.Exporter) gin.HandlerFunc {
	return func(c *gin.Context) {
		exp.ServeHTTP(c.Writer, c.Request)
	}
}
