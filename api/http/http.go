package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/eran-levy/tokenizer-gophercon/api/http/middleware"
	"github.com/eran-levy/tokenizer-gophercon/logger"
	"github.com/eran-levy/tokenizer-gophercon/service"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
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
type RestApiAdapter struct {
	cfg RestApiAdapterConfiguration
	srv *http.Server
	ts  service.TokenizerService
}

func New(cfg RestApiAdapterConfiguration, ts service.TokenizerService) *RestApiAdapter {
	if !cfg.IsDebugModeEnabled {
		gin.SetMode(gin.ReleaseMode)
	}
	return &RestApiAdapter{cfg: cfg, ts: ts}
}

func (s *RestApiAdapter) Start(fatalErrors chan<- error) {
	r := gin.New()
	r.Use(middleware.Logger())

	if s.cfg.IsDebugModeEnabled {
		//in case not using the default sever mux, register each one of the pprof routes
		pprof.Register(r)
	}
	r.GET("/health", health)
	r.GET("/readiness", readiness)
	//its possible also to activate middleware for a given group - just by pasing the middleware to r.Group
	v1 := r.Group("/v1")
	v1.POST("/tokenize", s.tokenizeTextHandler)

	srv := &http.Server{
		Addr:         s.cfg.HttpAddress,
		ReadTimeout:  s.cfg.ReadRequestTimeout,
		WriteTimeout: s.cfg.WriteResponseTimeout,
		Handler:      r,
	}
	s.srv = srv
	if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
		logger.Log.Errorf("http server closed: %s\n", err)
		fatalErrors <- err
	}
}

func (s *RestApiAdapter) Close(ctx context.Context) error {
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
	return nil
}
