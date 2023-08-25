package redis

import (
	"context"
	"errors"

	"go-boilerplate/config"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
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
		log.Panic().Err(ErrorRedisClientConnectionFailed).Msg(status.Val())
	}

	log.Info().Msgf("[CACHE] redis \t%s", status)
}
