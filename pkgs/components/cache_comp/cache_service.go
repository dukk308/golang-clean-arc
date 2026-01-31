package redis_component

import (
	"context"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

type CacheService struct {
	cache  *cache.Cache
	redis  redis.UniversalClient
	prefix string
}

type CacheItem struct {
	Ctx            context.Context
	Key            string
	Value          interface{}
	TTL            time.Duration
	Do             func(*cache.Item) (interface{}, error)
	SkipLocalCache bool
}

func NewCacheService(redisClient redis.UniversalClient, config *RedisCacheConfig) *CacheService {
	opts := &cache.Options{
		Redis: redisClient,
	}

	if config.EnableLocalCache {
		localTTL := time.Minute
		if config.LocalCacheTTL > 0 {
			localTTL = time.Duration(config.LocalCacheTTL) * time.Second
		}

		cacheSize := 1000
		if config.LocalCacheSize > 0 {
			cacheSize = config.LocalCacheSize
		}

		opts.LocalCache = cache.NewTinyLFU(cacheSize, localTTL)
	}

	prefix := config.Prefix
	if prefix != "" {
		prefix = prefix + ":"
	}

	return &CacheService{
		cache:  cache.New(opts),
		redis:  redisClient,
		prefix: prefix,
	}
}

func (c *CacheService) buildKey(key string) string {
	return c.prefix + key
}

func (c *CacheService) Set(item *CacheItem) error {
	return c.cache.Set(&cache.Item{
		Ctx:            item.Ctx,
		Key:            c.buildKey(item.Key),
		Value:          item.Value,
		TTL:            item.TTL,
		SkipLocalCache: item.SkipLocalCache,
	})
}

func (c *CacheService) Get(ctx context.Context, key string, value interface{}) error {
	return c.cache.Get(ctx, c.buildKey(key), value)
}

func (c *CacheService) Exists(ctx context.Context, key string) bool {
	return c.cache.Exists(ctx, c.buildKey(key))
}

func (c *CacheService) Delete(ctx context.Context, key string) error {
	return c.cache.Delete(ctx, c.buildKey(key))
}

func (c *CacheService) DeleteMultiple(ctx context.Context, keys ...string) error {
	prefixedKeys := make([]string, len(keys))
	for i, k := range keys {
		prefixedKeys[i] = c.buildKey(k)
	}

	for _, key := range prefixedKeys {
		if err := c.cache.Delete(ctx, key); err != nil {
			return err
		}
	}
	return nil
}

func (c *CacheService) Once(item *CacheItem) error {
	return c.cache.Once(&cache.Item{
		Ctx:            item.Ctx,
		Key:            c.buildKey(item.Key),
		Value:          item.Value,
		TTL:            item.TTL,
		Do:             item.Do,
		SkipLocalCache: item.SkipLocalCache,
	})
}

func (c *CacheService) SetPrefix(prefix string) {
	if prefix != "" {
		c.prefix = prefix + ":"
	} else {
		c.prefix = ""
	}
}

func (c *CacheService) GetRedisClient() redis.UniversalClient {
	return c.redis
}

func (c *CacheService) Close() error {
	if c.redis != nil {
		return c.redis.Close()
	}
	return nil
}

func (r *RedisClient) CacheSet(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   r.Key(key),
		Value: value,
		TTL:   ttl,
	})
}

func (r *RedisClient) CacheGet(ctx context.Context, key string, value interface{}) error {
	return r.cache.Get(ctx, r.Key(key), value)
}

func (r *RedisClient) CacheExists(ctx context.Context, key string) bool {
	return r.cache.Exists(ctx, r.Key(key))
}

func (r *RedisClient) CacheDelete(ctx context.Context, key string) error {
	return r.cache.Delete(ctx, r.Key(key))
}

func (r *RedisClient) CacheOnce(ctx context.Context, key string, value interface{}, ttl time.Duration, do func(*cache.Item) (interface{}, error)) error {
	return r.cache.Once(&cache.Item{
		Ctx:   ctx,
		Key:   r.Key(key),
		Value: value,
		TTL:   ttl,
		Do:    do,
	})
}

func (r *RedisClient) CacheSetSkipLocal(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.cache.Set(&cache.Item{
		Ctx:            ctx,
		Key:            r.Key(key),
		Value:          value,
		TTL:            ttl,
		SkipLocalCache: true,
	})
}

func (r *RedisClient) GetCache() *cache.Cache {
	return r.cache
}

func (r *RedisClient) ConfigureCache(localCacheSize int, localCacheTTL time.Duration) {
	r.cache = cache.New(&cache.Options{
		Redis:      r.c,
		LocalCache: cache.NewTinyLFU(localCacheSize, localCacheTTL),
	})
}

func (r *RedisClient) DisableLocalCache() {
	r.cache = cache.New(&cache.Options{
		Redis: r.c,
	})
}
