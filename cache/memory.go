package cache

import (
	"context"

	"github.com/lafriks/ttlcache/v3"
)

type memoryCache[T any] struct {
	cache *ttlcache.Cache[string, T]
}

func newMemoryCache[T any](opts ...CacheOption) (CacheInstance[T], error) {
	opt := newCacheOptions(opts...)
	c := ttlcache.New(ttlcache.WithTTL[string, T](opt.TTL))
	return &memoryCache[T]{c}, nil
}

func (c *memoryCache[T]) Get(ctx context.Context, key string, opts ...ItemOption[T]) (T, error) {
	i, err := c.cache.Get(key)
	var val T
	if err != nil {
		return val, err
	}
	if i.IsExpired() {
		return val, nil
	}
	return i.Value(), nil
}

func (c *memoryCache[T]) Set(ctx context.Context, key string, value T, opts ...ItemOption[T]) error {
	opt := newItemOptions(opts...)
	ttl := opt.TTL
	if ttl == 0 {
		ttl = ttlcache.DefaultTTL
	}
	_ = c.cache.Set(key, value, ttl)
	return nil
}

func (c *memoryCache[T]) Delete(ctx context.Context, key string) error {
	c.cache.Delete(key)
	return nil
}

func (c *memoryCache[T]) Close() {
	c.cache.DeleteAll()
	c.cache = nil
}