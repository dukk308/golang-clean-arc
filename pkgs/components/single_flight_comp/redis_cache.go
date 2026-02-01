package single_flight_comp

import (
	"context"
	"time"

	"github.com/go-redis/cache/v9"
)

type RedisCacheSingleflight struct {
	cache *cache.Cache
	group Singleflight
	keyFn func(string) string
}

func NewRedisCacheSingleflight(cache *cache.Cache, group Singleflight, keyFn func(string) string) *RedisCacheSingleflight {
	if keyFn == nil {
		keyFn = func(k string) string { return k }
	}
	return &RedisCacheSingleflight{cache: cache, group: group, keyFn: keyFn}
}

func (r *RedisCacheSingleflight) CacheOnce(
	ctx context.Context,
	key string,
	value interface{},
	ttl time.Duration,
	fetch func(*cache.Item) (interface{}, error),
) (*Result, error) {
	fullKey := r.keyFn(key)
	result, err, shared := r.group.Do(fullKey, func() (interface{}, error) {
		err := r.cache.Once(&cache.Item{
			Ctx:   ctx,
			Key:   fullKey,
			Value: value,
			TTL:   ttl,
			Do:    fetch,
		})
		return value, err
	})
	return &Result{Value: result, Err: err, Shared: shared}, err
}

func (r *RedisCacheSingleflight) CacheGet(
	ctx context.Context,
	key string,
	value interface{},
	ttl time.Duration,
	fetch func() (interface{}, error),
) (*Result, error) {
	fullKey := r.keyFn(key)
	result, err, shared := r.group.Do(fullKey, func() (interface{}, error) {
		err := r.cache.Get(ctx, fullKey, value)
		if err == nil {
			return value, nil
		}
		fetchedValue, fetchErr := fetch()
		if fetchErr != nil {
			return nil, fetchErr
		}
		_ = r.cache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   fullKey,
			Value: fetchedValue,
			TTL:   ttl,
		})
		return fetchedValue, nil
	})
	return &Result{Value: result, Err: err, Shared: shared}, err
}

func (r *RedisCacheSingleflight) Forget(key string) {
	r.group.Forget(r.keyFn(key))
}
