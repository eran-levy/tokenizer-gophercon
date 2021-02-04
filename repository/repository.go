package repository

import (
	"context"
	"github.com/eran-levy/tokenizer-gophercon/repository/model"
	"time"
)

type Persistence interface {
	StoreMetadata(ctx context.Context, mtd model.TokenizeTextMetadata) error
	Close() error
	IsServiceHealthy(ctx context.Context) (bool, error)
}

type Config struct {
	Dsn                   string
	ConnectionMaxLifetime time.Duration
	MaxOpenConnections    int
	MaxIdleConnections    int
}
