package redis_component

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/go-redsync/redsync/v4"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	prefix  string
	c       redis.UniversalClient
	timeout time.Duration
	cache   *cache.Cache
	redsync *redsync.Redsync
}

type CacheKeyBuilder struct {
	prefix string
}

func NewCacheKeyBuilder(prefix string) *CacheKeyBuilder {
	return &CacheKeyBuilder{prefix: prefix}
}

func NewStandAloneRedisClient(prefix string, timeout time.Duration, c *redis.Client) ICacheService {
	if prefix != "" {
		prefix = prefix + ":"
	}

	cacheInstance := cache.New(&cache.Options{
		Redis:      c,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	client := RedisClient{
		prefix:  prefix,
		c:       c,
		timeout: timeout,
		cache:   cacheInstance,
	}
	return &client
}

func (r *RedisClient) Key(key string) string {
	return r.prefix + key
}

func (r *RedisClient) SetPrefix(prefix string) {
	if prefix != "" {
		r.prefix = prefix + ":"
	} else {
		r.prefix = ""
	}
}

func (r *RedisClient) MustGet(context context.Context, key string) (*string, error) {
	val, err := r.c.Get(context, r.Key(key)).Result()
	if err != nil {
		return nil, err
	}

	return &val, nil
}

func (r *RedisClient) Get(ctx context.Context, key string) (*string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	val, err := r.c.Get(ctx, r.Key(key)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	return &val, nil
}

func (r *RedisClient) marshal(value interface{}) (string, error) {
	var strValue string

	switch v := value.(type) {
	case string:
		strValue = v
	case []byte:
		strValue = string(v)
	default:
		data, err := json.Marshal(value)
		if err != nil {
			return "", err
		}
		strValue = string(data)
	}

	return strValue, nil
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	strValue, err := r.marshal(value)
	if err != nil {
		return err
	}

	return r.c.Set(ctx, r.Key(key), strValue, 0).Err()
}

func (r *RedisClient) Delete(ctx context.Context, key ...string) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	keys := make([]string, len(key))
	for i, k := range key {
		keys[i] = r.Key(k)
	}
	return r.c.Del(ctx, keys...).Err()
}

func (r *RedisClient) Close() error {
	return r.c.Close()
}

func (r *RedisClient) SetEx(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	strValue, err := r.marshal(value)
	if err != nil {
		return err
	}
	return r.c.Set(ctx, r.Key(key), strValue, expiration).Err()
}

func (r *RedisClient) IsExist(ctx context.Context, key string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	exists, err := r.c.Exists(ctx, r.Key(key)).Result()
	return exists == 1, err
}

func (r *RedisClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	ttl, err := r.c.TTL(ctx, r.Key(key)).Result()
	return ttl, err
}

func (r *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	return r.c.Expire(ctx, r.Key(key), expiration).Err()
}

func (r *RedisClient) Extends(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	remain, err := r.TTL(ctx, key)
	if err != nil {
		return false, err
	}

	if remain <= 0 {
		return false, nil
	}

	newExpire := remain + expiration
	ok, err := r.c.Expire(ctx, r.Key(key), newExpire).Result()
	return ok, err
}

func (r *RedisClient) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.c.Ping(ctx).Err()
}

func (r *RedisClient) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.c.Scan(ctx, cursor, match, count).Result()
}

func (r *RedisClient) SAdd(ctx context.Context, key string, members ...interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.c.SAdd(ctx, r.Key(key), members...).Result()
}

func (r *RedisClient) SRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.c.SRem(ctx, r.Key(key), members...).Result()
}

func (r *RedisClient) SMembers(ctx context.Context, key string) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.c.SMembers(ctx, r.Key(key)).Result()
}

func (r *RedisClient) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.c.SIsMember(ctx, r.Key(key), member).Result()
}

func (r *RedisClient) SPop(ctx context.Context, key string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.c.SPop(ctx, r.Key(key)).Result()
}

func (r *RedisClient) SCard(ctx context.Context, key string) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.c.SCard(ctx, r.Key(key)).Result()
}

func (r *RedisClient) GetInt(ctx context.Context, key string) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.c.Get(ctx, r.Key(key)).Int64()
}

func (r *RedisClient) Incr(ctx context.Context, key string) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.c.Incr(ctx, r.Key(key)).Result()
}

func (r *RedisClient) Decr(ctx context.Context, key string) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.c.Decr(ctx, r.Key(key)).Result()
}

func (r *RedisClient) Publish(ctx context.Context, channel string, message interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.c.Publish(ctx, r.Key(channel), message).Result()
}

func (r *RedisClient) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var prefixed []string
	for _, ch := range channels {
		prefixed = append(prefixed, r.Key(ch))
	}
	return r.c.Subscribe(ctx, prefixed...)
}

func (r *RedisClient) GetWithPrefix(ctx context.Context, prefix, key string) (*string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	val, err := r.c.Get(ctx, prefix+":"+key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	return &val, nil
}

func (r *RedisClient) SGetWithPrefix(ctx context.Context, prefix, key string) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	members, err := r.c.SMembers(ctx, prefix+":"+key).Result()
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (r *RedisClient) C() redis.UniversalClient {
	return r.c
}

func (r *RedisClient) HSet(ctx context.Context, key string, values ...interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.c.HSet(ctx, r.Key(key), values...).Err()
}

func (r *RedisClient) HGet(ctx context.Context, key string, field string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.c.HGet(ctx, r.Key(key), field).Result()
}

func (r *RedisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.c.HGetAll(ctx, r.Key(key)).Result()
}

func (r *RedisClient) HMGet(ctx context.Context, key string, fields ...string) ([]interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.c.HMGet(ctx, r.Key(key), fields...).Result()
}

func (r *RedisClient) HIncrByFloat(ctx context.Context, key string, field string, incr float64) (float64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.c.HIncrByFloat(ctx, r.Key(key), field, incr).Result()
}
