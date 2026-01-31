package redis_component

import (
	"time"

	"github.com/dukk308/golang-clean-arch-starter/pkgs/logger"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type RedisComponent struct {
	id           string
	config       *RedisConfig
	addrs        []string
	pass         string
	prefix       string
	client       ICacheService
	cacheService *CacheService
	timeout      time.Duration
}

func NewRedisComponent(config *RedisConfig, logger logger.Logger) *RedisComponent {
	component := &RedisComponent{
		config: config,
	}

	component.addrs = config.Addrs
	component.pass = config.Password
	component.prefix = config.Prefix
	component.timeout = time.Duration(config.Timeout) * time.Millisecond

	if err := component.Activate(logger); err != nil {
		logger.Errorf("Failed to activate Redis component: %v", err)
	}

	return component
}

func (r *RedisComponent) GetClient() ICacheService {
	return r.client
}

func (r *RedisComponent) GetCacheService() *CacheService {
	return r.cacheService
}

func (r *RedisComponent) InitializeCacheService(cacheConfig *RedisCacheConfig) {
	if r.client != nil {
		r.cacheService = NewCacheService(r.client.C(), cacheConfig)
	}
}

func (r *RedisComponent) Activate(logger logger.Logger) error {
	if err := r.connect(logger); err != nil {
		return err
	}

	return nil
}

func (r *RedisComponent) ID() string {
	return r.id
}

func (r *RedisComponent) Stop() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

func (r *RedisComponent) connect(logger logger.Logger) error {
	logger.Info("Connecting to Redis ", r.addrs)

	logger.Debugf("Redis address: %v", r.addrs[0])
	logger.Debugf("Redis password: %v", r.pass)
	logger.Debugf("Redis prefix: %v", r.prefix)
	logger.Debugf("Redis pool size: %v", r.config.PoolSize)
	logger.Debugf("Redis min idle conns: %v", r.config.MinIdleConns)
	logger.Debugf("Redis max retries: %v", r.config.MaxRetries)
	logger.Debugf("Redis dial timeout: %v", r.config.DialTimeout)
	logger.Debugf("Redis read timeout: %v", r.config.ReadTimeout)

	opts := &redis.Options{
		Addr:       r.addrs[0],
		Password:   r.pass,
		DB:         0,
		MaxRetries: 3,
	}

	r.applyStandaloneOptions(opts)

	client := redis.NewClient(opts)
	r.withInstrumentation(client, logger)
	r.client = NewStandAloneRedisClient(r.prefix, r.timeout, client)

	logger.Info("Successfully connected to Redis")
	return nil
}

func (r *RedisComponent) applyStandaloneOptions(opts *redis.Options) {
	if r.config == nil {
		return
	}

	if r.config.PoolSize > 0 {
		opts.PoolSize = r.config.PoolSize
	}
	if r.config.MinIdleConns > 0 {
		opts.MinIdleConns = r.config.MinIdleConns
	}
	if r.config.MaxRetries > 0 {
		opts.MaxRetries = r.config.MaxRetries
	}
	if r.config.DialTimeout > 0 {
		opts.DialTimeout = time.Duration(r.config.DialTimeout) * time.Second
	}
	if r.config.ReadTimeout > 0 {
		opts.ReadTimeout = time.Duration(r.config.ReadTimeout) * time.Second
	}
	if r.config.WriteTimeout > 0 {
		opts.WriteTimeout = time.Duration(r.config.WriteTimeout) * time.Second
	}
	if r.config.PoolTimeout > 0 {
		opts.PoolTimeout = time.Duration(r.config.PoolTimeout) * time.Second
	}
}

func (r *RedisComponent) withInstrumentation(client redis.UniversalClient, logger logger.Logger) {
	enableTracing := true
	if r.config != nil {
		enableTracing = r.config.EnableTracing
	}

	if enableTracing {
		if err := redisotel.InstrumentTracing(client); err != nil {
			logger.Warnf("Failed to instrument Redis client with tracing: %v", err)
		} else {
			logger.Info("Redis client instrumented with OpenTelemetry tracing")
		}
	}
}
