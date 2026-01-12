package domain

import (
	"context"
	"time"

	"github.com/ali-nur31/mile-do/internal/repository/db"
)

type CreateRefreshTokenInput struct {
	UserID    int32
	TokenHash string
	ExpiresAt time.Time
}

type RefreshTokenOutput struct {
	ID        int64
	UserID    int32
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
}

func ToRefreshTokenOutput(token *repo.RefreshToken) *RefreshTokenOutput {
	return &RefreshTokenOutput{
		ID:        token.ID,
		UserID:    token.UserID,
		TokenHash: token.TokenHash,
		ExpiresAt: token.ExpiresAt.Time,
		CreatedAt: token.CreatedAt.Time,
	}
}

type RefreshTokenService interface {
	GetRefreshTokenByUserID(ctx context.Context, qtx repo.Querier, userId int32) (*RefreshTokenOutput, error)
	CreateRefreshToken(ctx context.Context, input CreateRefreshTokenInput) error
	DeleteRefreshTokenByUserID(ctx context.Context, userId int32) error
}
