package redis_component

import (
	"flag"
	"strings"
)

var (
	redisAddrs   = flag.String("redis-addrs", "localhost:6379", "Redis addresses (comma-separated)")
	redisPass    = flag.String("redis-pass", "", "Redis password")
	redisTimeout = flag.Int("redis-timeout", 20000, "Redis timeout in miliseconds")
)

func LoadRedisConfig(serviceName string) *RedisConfig {
	addrs := strings.Split(*redisAddrs, ",")
	for i := range addrs {
		addrs[i] = strings.TrimSpace(addrs[i])
	}
	prefix := strings.ToUpper(serviceName)

	return &RedisConfig{
		Addrs:         addrs,
		Password:      *redisPass,
		Prefix:        prefix,
		Timeout:       *redisTimeout,
		EnableTracing: true,
	}
}

