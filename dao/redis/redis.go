package redis

import (
	"demo_webapp/setting"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
)

func Init(cfg *setting.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port),
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})
	return
}

func Close() {
	defer client.Close()
}
