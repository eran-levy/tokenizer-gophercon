package http

import (
	"github.com/eran-levy/tokenizer-gophercon/logger"
	"github.com/eran-levy/tokenizer-gophercon/telemetry"
	"github.com/gin-gonic/gin"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		//before request
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		c.Next()
		//after request
		d := time.Since(start).Milliseconds()
		logger.Log.With("path", path).With("query_params", raw).With("duration", d).Info("API request")
		telemetry.RecordAPIRequestLatencyValue(c.Request.Context(), d, path, telemetry.SuccessStatusValue)
	}
}
