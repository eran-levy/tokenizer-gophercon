package cache

import (
	"context"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key string, value []byte) error
	Get(ctx context.Context, key string) ([]byte, bool)
	Close() error
}

type Config struct {
	CacheSize      int
	CacheAddress   string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	ExpirationTime time.Duration
	MaxRetries     int
	//in case you would like to leave some idle conn
	//to save time on establishing new conns
	MinIdleConns int
}
