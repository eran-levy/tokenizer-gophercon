package local

import (
	"context"
	"github.com/eran-levy/tokenizer-gophercon/cache"
	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
)

type localCache struct {
	c *lru.Cache
}

func (l localCache) Close() error {
	return nil
}

func New(cfg cache.Config) (cache.Cache, error) {
	c, err := lru.New(cfg.CacheSize)
	if err != nil {
		return localCache{}, errors.Wrap(err, "init local cache failed")
	}
	return localCache{c: c}, nil
}

func (l localCache) Set(ctx context.Context, key string, value []byte) error {
	l.c.Add(key, value)
	return nil
}

func (l localCache) Get(ctx context.Context, key string) ([]byte, bool) {
	v, found := l.c.Get(key)
	if found {
		return v.([]byte), found
	}
	return []byte{}, found
}

func (l localCache) IsServiceHealthy(ctx context.Context) (bool, error) {
	return true, nil
}
