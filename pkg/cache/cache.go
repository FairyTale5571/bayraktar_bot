package cache

import (
	"time"

	"github.com/fairytale5571/bayraktar_bot/pkg/cache/memory_cache"
	"github.com/fairytale5571/bayraktar_bot/pkg/cache/redis_cache"
	errs "github.com/fairytale5571/bayraktar_bot/pkg/errorUtils"
	"github.com/fairytale5571/bayraktar_bot/pkg/logger"
	"github.com/fairytale5571/bayraktar_bot/pkg/storage"
	"github.com/fairytale5571/bayraktar_bot/pkg/storage/redis"
)

const (
	ttlMemory = 15 * time.Minute
	ttlRedis  = 30 * time.Minute
)

type Cache interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string)
}

type Config struct {
	rdb        *redis.Redis
	memory     *memory_cache.Memory
	redisCache *redis_cache.RedisCache
	logger     *logger.Wrapper

	stop chan struct{}
}

func SetupCache(rdb *redis.Redis) *Config {
	return &Config{
		rdb:        rdb,
		memory:     memory_cache.New(ttlMemory),
		redisCache: redis_cache.New(rdb),
		logger:     logger.New("cache"),
	}
}

func (c *Config) Get(key string) (string, error) {
	if res, err := c.memory.Get(key); err == nil {
		return res, nil
	}
	if res, err := c.rdb.Get(key, storage.Cache); err == nil {
		return res, nil
	}
	return "", errs.ErrorNotCached
}

func (c *Config) Set(key, value string) error {
	if err := c.memory.Set(key, value); err != nil {
		c.logger.Errorf("cant cache key %s in memory: %s", key, err)
		return errs.ErrorCantCacheMemory
	}
	if err := c.redisCache.Set(key, value, ttlRedis); err != nil {
		c.logger.Errorf("cant cache key %s in redis: %s", key, err)
		return errs.ErrorCantCacheRedis
	}
	return nil
}
