package redis_component

import (
	"context"
	"time"

	"github.com/go-redis/cache/v9"
	"golang.org/x/sync/singleflight"
)

func (r *RedisClient) SingleflightDo(key string, fn func() (interface{}, error)) (interface{}, error, bool) {
	result, err, shared := r.group.Do(r.Key(key), fn)
	return result, err, shared
}

func (r *RedisClient) SingleflightCacheOnce(
	ctx context.Context,
	key string,
	value interface{},
	ttl time.Duration,
	fetch func(*cache.Item) (interface{}, error),
) (*SingleflightResult, error) {
	fullKey := r.Key(key)
	
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
	
	return &SingleflightResult{
		Value:  result,
		Err:    err,
		Shared: shared,
	}, err
}

func (r *RedisClient) SingleflightCacheGet(
	ctx context.Context,
	key string,
	value interface{},
	ttl time.Duration,
	fetch func() (interface{}, error),
) (*SingleflightResult, error) {
	fullKey := r.Key(key)
	
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
	
	return &SingleflightResult{
		Value:  result,
		Err:    err,
		Shared: shared,
	}, err
}

func (r *RedisClient) SingleflightForget(key string) {
	r.group.Forget(r.Key(key))
}

func (r *RedisClient) SingleflightDoChan(key string, fn func() (interface{}, error)) <-chan singleflight.Result {
	return r.group.DoChan(r.Key(key), fn)
}
