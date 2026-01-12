package redis

import (
	"context"
	"time"

	"github.com/ali-nur31/mile-do/internal/domain"
	"github.com/redis/go-redis/v9"
)

type authRedisRepo struct {
	client *redis.Client
}

func NewAuthRedisRepo(client *redis.Client) domain.AuthCacheRepo {
	return &authRedisRepo{
		client: client,
	}
}

func (r *authRedisRepo) BlockToken(ctx context.Context, tokenID string, duration time.Duration) error {
	key := "blacklist:accessToken:" + tokenID
	return r.client.Set(ctx, key, "true", duration).Err()
}

func (r *authRedisRepo) IsTokenBlocked(ctx context.Context, tokenID string) (bool, error) {
	key := "blacklist:accessToken:" + tokenID
	exists, err := r.client.Exists(ctx, key).Result()
	return exists > 0, err
}
