package main

import (
	"context"
	"github.com/eran-levy/tokenizer-gophercon/api/http"
	"github.com/eran-levy/tokenizer-gophercon/cache"
	"github.com/eran-levy/tokenizer-gophercon/cache/redis"
	"github.com/eran-levy/tokenizer-gophercon/config"
	"github.com/eran-levy/tokenizer-gophercon/logger"
	"github.com/eran-levy/tokenizer-gophercon/service"
	"github.com/eran-levy/tokenizer-gophercon/telemetry"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	logger.New(logger.Config{LogLevel: cfg.Service.LogLevel, ApplicationId: cfg.Service.AppId})
	defer logger.Close()
	telem, flush, err := telemetry.New(telemetry.Config{ApplicationID: cfg.Service.AppId, ServiceName: cfg.Service.AppId, AgentEndpoint: cfg.Telemetry.TracingAgentEndpoint})
	if err != nil {
		logger.Log.Fatal(err)
	}
	defer flush()
	//c, err := local.New(cache.Config{CacheSize: cfg.Cache.CacheSize})
	c, err := redis.New(cache.Config{CacheAddress: cfg.Cache.CacheAddress, ReadTimeout: cfg.Cache.ReadTimeout,
		ExpirationTime: cfg.Cache.ExpirationTime})
	if err != nil {
		logger.Log.Fatal(err)
	}
	fatalErrors := make(chan error, 1)
	ts := service.New(c)
	srv := http.New(http.RestApiAdapterConfiguration{HttpAddress: cfg.RESTApiAdapter.HttpAddress,
		TerminationTimeout:   cfg.RESTApiAdapter.TerminationTimeout,
		ReadRequestTimeout:   cfg.RESTApiAdapter.ReadRequestTimeout,
		WriteResponseTimeout: cfg.RESTApiAdapter.WriteResponseTimeout,
		IsDebugModeEnabled:   cfg.Service.DebugModeEnabled}, ts, telem)
	go srv.Start(fatalErrors)

	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, os.Interrupt, syscall.SIGTERM)
	select {
	case err := <-fatalErrors:
		// unexepected failures should arrive in here
		logger.Log.Errorf("fatal error signal received %s\n", err)
		os.Exit(1)
	case sig := <-gracefulShutdown:
		// try to terminal all gracefully here
		logger.Log.Infof("observed terminal signal %v\n", sig)
		err := srv.Close(context.Background())
		if err != nil {
			logger.Log.Errorf("could not close servers gracefully %s\n", err)
		}
		err = c.Close()
		if err != nil {
			logger.Log.Errorf("could not close cache %s \n ", err)
		}
	}

}
