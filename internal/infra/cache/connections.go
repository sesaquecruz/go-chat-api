package cache

import (
	"fmt"

	"github.com/sesaquecruz/go-chat-api/config"

	"github.com/redis/go-redis/v9"
)

func RedisConnection(cfg *config.CacheConfig) *redis.Client {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	return rdb
}
