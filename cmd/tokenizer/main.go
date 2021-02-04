package main

import (
	"context"
	"github.com/eran-levy/tokenizer-gophercon/api/grpc"
	"github.com/eran-levy/tokenizer-gophercon/api/http"
	"github.com/eran-levy/tokenizer-gophercon/cache"
	"github.com/eran-levy/tokenizer-gophercon/cache/local"
	"github.com/eran-levy/tokenizer-gophercon/cache/redis"
	"github.com/eran-levy/tokenizer-gophercon/config"
	"github.com/eran-levy/tokenizer-gophercon/logger"
	"github.com/eran-levy/tokenizer-gophercon/repository"
	"github.com/eran-levy/tokenizer-gophercon/repository/mysql"
	"github.com/eran-levy/tokenizer-gophercon/service"
	"github.com/eran-levy/tokenizer-gophercon/telemetry"
	"log"
	nhtp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	logger.New(logger.Config{LogLevel: cfg.Service.LogLevel, ApplicationId: cfg.Service.AppId})
	defer logger.Close()
	//setup telemetry
	telem, flush, err := telemetry.New(telemetry.Config{ApplicationID: cfg.Service.AppId, ServiceName: cfg.Service.AppId, AgentEndpoint: cfg.Telemetry.TracingAgentEndpoint})
	if err != nil {
		logger.Log.Fatal(err)
	}
	defer flush()
	//errors channel to gracefully close resources
	fatalErrors := make(chan error, 1)

	//c, err := local.New(cache.Config{CacheSize: cfg.Cache.CacheSize})
	var c cache.Cache
	c, err = redis.New(cache.Config{CacheAddress: cfg.Cache.CacheAddress, ReadTimeout: cfg.Cache.ReadTimeout,
		ExpirationTime: cfg.Cache.ExpirationTime})
	if err != nil {
		logger.Log.Errorf("could not connect to distributed cache, init local %s", err)
		c, err = local.New(cache.Config{CacheSize: cfg.Cache.CacheSize})
		if err != nil {
			logger.Log.Fatal(err)
		}
	}
	//setup persistence
	repo, err := mysql.New(repository.Config{User: cfg.Database.User, Passwd: cfg.Database.Passwd,
		Address: cfg.Database.Address, DBName: cfg.Database.DBName, ConnectionMaxLifetime: cfg.Database.ConnectionMaxLifetime,
		MaxOpenConnections: cfg.Database.MaxOpenConnections, MaxIdleConnections: cfg.Database.MaxIdleConnections}, telem)
	if err != nil {
		logger.Log.Fatal(err)
	}
	htClient := &nhtp.Client{Timeout: time.Second * 30}
	ts := service.New(c, repo, telem, htClient)
	srv := http.New(http.RestApiAdapterConfiguration{HttpAddress: cfg.RESTApiAdapter.HttpAddress,
		TerminationTimeout:   cfg.RESTApiAdapter.TerminationTimeout,
		ReadRequestTimeout:   cfg.RESTApiAdapter.ReadRequestTimeout,
		WriteResponseTimeout: cfg.RESTApiAdapter.WriteResponseTimeout,
		IsDebugModeEnabled:   cfg.Service.DebugModeEnabled}, ts, telem)
	go srv.Start(fatalErrors)

	gSrv := grpc.New(grpc.Config{GrpcAddress: cfg.GRPCApiAdapter.GrpcAddress,
		MaxConnectionAge:      cfg.GRPCApiAdapter.MaxConnectionAge,
		MaxConnectionAgeGrace: cfg.GRPCApiAdapter.MaxConnectionAgeGrace}, ts, telem)
	go gSrv.Start(fatalErrors)

	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, os.Interrupt, syscall.SIGTERM)
	select {
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
		err = repo.Close()
		if err != nil {
			logger.Log.Errorf("could not close repository handlers %s \n ", err)
		}
		gSrv.Close()
	case err := <-fatalErrors:
		// unexepected failures should arrive in here
		logger.Log.Errorf("fatal error signal received %s\n", err)
		os.Exit(1)
	}

}
