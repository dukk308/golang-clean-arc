package redis_component

import (
	"context"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

type ICacheService interface {
	Get(context context.Context, key string) (*string, error)
	MustGet(context context.Context, key string) (*string, error)
	Set(context context.Context, key string, value interface{}) error
	Delete(context context.Context, key ...string) error
	Close() error
	SetEx(context context.Context, key string, value interface{}, expiration time.Duration) error
	IsExist(context context.Context, key string) (bool, error)
	Extends(context context.Context, key string, expiration time.Duration) (bool, error)
	Ping(context context.Context) error
	Scan(context context.Context, cursor uint64, match string, count int64) ([]string, uint64, error)
	SetPrefix(prefix string)
	SAdd(ctx context.Context, key string, members ...interface{}) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	TTL(ctx context.Context, key string) (time.Duration, error)
	SRem(ctx context.Context, key string, members ...interface{}) (int64, error)
	SMembers(ctx context.Context, key string) ([]string, error)
	SIsMember(ctx context.Context, key string, member interface{}) (bool, error)
	SPop(ctx context.Context, key string) (string, error)
	SCard(ctx context.Context, key string) (int64, error)
	Publish(ctx context.Context, channel string, message interface{}) (int64, error)
	Subscribe(ctx context.Context, channels ...string) *redis.PubSub
	GetWithPrefix(context context.Context, prefix, key string) (*string, error)
	SGetWithPrefix(ctx context.Context, prefix, key string) ([]string, error)
	Incr(ctx context.Context, key string) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)
	C() redis.UniversalClient
	GetInt(ctx context.Context, key string) (int64, error)
	HSet(ctx context.Context, key string, values ...interface{}) error
	HGet(ctx context.Context, key string, field string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HMGet(ctx context.Context, key string, fields ...string) ([]interface{}, error)
	HIncrByFloat(ctx context.Context, key string, field string, incr float64) (float64, error)

	CacheSet(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	CacheGet(ctx context.Context, key string, value interface{}) error
	CacheExists(ctx context.Context, key string) bool
	CacheDelete(ctx context.Context, key string) error
	CacheOnce(ctx context.Context, key string, value interface{}, ttl time.Duration, do func(*cache.Item) (interface{}, error)) error
	CacheSetSkipLocal(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	GetCache() *cache.Cache
	ConfigureCache(localCacheSize int, localCacheTTL time.Duration)
	DisableLocalCache()
}
