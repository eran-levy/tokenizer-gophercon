package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/eran-levy/tokenizer-gophercon/logger"
	"github.com/eran-levy/tokenizer-gophercon/service"
	"github.com/eran-levy/tokenizer-gophercon/telemetry"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"net/http"
	"time"
)

var (
	requestNotValid = errors.New("API request is not valid")
)

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type RestApiAdapterConfiguration struct {
	HttpAddress          string
	TerminationTimeout   time.Duration
	ReadRequestTimeout   time.Duration
	WriteResponseTimeout time.Duration
	IsDebugModeEnabled   bool
}
type restApiAdapter struct {
	cfg       RestApiAdapterConfiguration
	srv       *http.Server
	ts        service.TokenizerService
	telemetry telemetry.Telemetry
}

func New(cfg RestApiAdapterConfiguration, ts service.TokenizerService, telemetry telemetry.Telemetry) *restApiAdapter {
	if !cfg.IsDebugModeEnabled {
		gin.SetMode(gin.ReleaseMode)
	}
	return &restApiAdapter{cfg: cfg, ts: ts, telemetry: telemetry}
}

func (s *restApiAdapter) Start(fatalErrors chan<- error) {
	r := gin.New()
	r.Use(Logger())
	r.Use(otelgin.Middleware(s.telemetry.Config.ServiceName))

	if s.cfg.IsDebugModeEnabled {
		//in case not using the default sever mux, register each one of the pprof routes
		pprof.Register(r)
	}
	//TODO: consider moving to telemetry pkg
	exp, err := telemetry.GetMeterHandlerToServe()
	if err != nil {
		logger.Log.Errorf("could not init metrics exporter: %s\n", err)
		fatalErrors <- err
	}
	r.GET("/metrics", MetricHandler(exp))
	r.GET("/health", s.health)
	r.GET("/readiness", s.readiness)
	// its also possible to set timeout for specific reoute
	r.GET("/demo", timeout.New(timeout.WithTimeout(10*time.Second), timeout.WithHandler(demo)))
	//its possible also to activate middleware for a given group - just by pasing the middleware to r.Group
	v1 := r.Group("/v1")
	v1.POST("/tokenize", s.tokenizeTextHandler)

	srv := &http.Server{
		Addr:         s.cfg.HttpAddress,
		ReadTimeout:  s.cfg.ReadRequestTimeout,
		WriteTimeout: s.cfg.WriteResponseTimeout,
		Handler:      http.TimeoutHandler(r, 2*time.Second, ""),
	}
	s.srv = srv
	if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
		logger.Log.Infof("http server closed %s", err)
		fatalErrors <- err
	}
}

func (s *restApiAdapter) Close(ctx context.Context) error {
	if s.srv == nil {
		return fmt.Errorf("could not shutdown http server cause it was not set in the adapter")
	}
	// The context is used to inform the server it has the termination timeout to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(ctx, s.cfg.TerminationTimeout)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}
	logger.Log.Info("gracefully closed http server")
	return nil
}
