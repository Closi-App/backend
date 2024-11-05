package redis

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func NewRedis(cfg *viper.Viper) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.GetString("redis.host"), cfg.GetInt("redis.port")),
		Password: cfg.GetString("redis.password"),
		DB:       cfg.GetInt("redis.database"),
	})
}
