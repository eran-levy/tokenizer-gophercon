package local

import (
	"github.com/eran-levy/tokenizer-gophercon/cache"
	"github.com/stretchr/testify/assert"
	"testing"
)

// for demonstration purposes
func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		cfg      cache.Config
		isFailed bool
		lc       localCache
	}{
		{
			name:     "fail not provided cache size",
			cfg:      cache.Config{},
			isFailed: true,
			lc:       localCache{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := New(tt.cfg)
			assert.NotNil(t, err)
			assert.Equal(t, tt.lc, c)
		})
	}
}

