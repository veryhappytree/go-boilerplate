package redis

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"go-boilerplate/config"

	"github.com/go-redis/redis/v9"
)

const successStatus = "PONG"

var Client *redis.Client
var ErrorRedisClientConnectionFailed = errors.New("[CACHE] redis client connection was failed")

func Setup(ctx context.Context, cfg config.RedisConfig) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost,
		Password: "",
		DB:       0,
		PoolSize: 1,
	})

	Client = rdb

	status := Client.Ping(ctx)

	if status.Val() != successStatus {
		slog.Error("[CACHE]", "error", ErrorRedisClientConnectionFailed)
	}

	slog.Info("[CACHE]", "message", fmt.Sprintf("redis status %s", status))
}
