package redis

import (
	"fmt"

	"github.com/go-redis/redis/v8"

	"github.com/KarolosLykos/sportsnews/config"
)

// New creates new connection to reds.
func New(cfg *config.Config) *redis.Client {
	addr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)

	return redis.NewClient(&redis.Options{Addr: addr})
}
