package redis_db

import (
	"context"
	"log/slog"

	"github.com/ali-nur31/mile-do/config"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Rdb *redis.Client
}

func InitializeRedisConnection(ctx context.Context, cfg *config.Redis) (*Redis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		Protocol: cfg.Protocol,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		slog.Error("failed to connect to Redis", "error", err)
		return nil, err
	}

	return &Redis{Rdb: rdb}, nil
}
