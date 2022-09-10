package redis_cache

import (
	"time"

	"github.com/fairytale5571/bayraktar_bot/pkg/storage"
	"github.com/fairytale5571/bayraktar_bot/pkg/storage/redis"
)

type RedisCache struct {
	rdb *redis.Redis
}

func New(rdb *redis.Redis) *RedisCache {
	return &RedisCache{
		rdb: rdb,
	}
}

func (c *RedisCache) Set(key, value string, ttl time.Duration) error {
	if err := c.rdb.SetTTL(key, value, storage.Cache, ttl); err != nil {
		return err
	}
	return nil
}

func (c *RedisCache) Get(key string) (string, error) {
	res, err := c.rdb.Get(key, storage.Cache)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (c *RedisCache) Delete(key string) {
	c.rdb.Delete(key, storage.Cache)
}

func (c *RedisCache) Cleanup() {
}
