package cache

import (
	"context"
	"time"

	"github.com/nullexp/limiter-x/internal/port/driven"
	cache "github.com/patrickmn/go-cache"
)

type MemoryClient struct {
	client *cache.Cache

	defaultExpiration, cleanupInterval time.Duration
}

func NewMemoryClient(defaultExpiration, cleanupInterval time.Duration) driven.Cache {
	return &MemoryClient{defaultExpiration: defaultExpiration, cleanupInterval: cleanupInterval}
}

func (rc *MemoryClient) Connect() error {
	rc.client = cache.New(rc.defaultExpiration, rc.cleanupInterval)

	return nil
}

func (rc *MemoryClient) Disconnect() error {
	rc.client.Flush()

	return nil
}

func (rc *MemoryClient) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	rc.client.Set(key, value, expiration)

	return nil
}

func (rc *MemoryClient) Fetch(ctx context.Context, key string) ([]byte, error) {
	value, exist := rc.client.Get(key)

	if !exist {
		return nil, driven.ErrCacheMissed
	}

	return value.([]byte), nil
}

func (rc *MemoryClient) Delete(ctx context.Context, key string) error {
	rc.client.Delete(key)

	return nil
}
