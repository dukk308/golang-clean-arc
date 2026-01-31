package redis_component

type RedisConfig struct {
	Addrs         []string
	Password      string
	Prefix        string
	PoolSize      int
	MinIdleConns  int
	MaxRetries    int
	DialTimeout   int
	ReadTimeout   int
	WriteTimeout  int
	PoolTimeout   int
	Timeout       int
	EnableTracing bool
}

type RedisCacheConfig struct {
	Prefix           string
	EnableLocalCache bool
	LocalCacheSize   int
	LocalCacheTTL    int
}
