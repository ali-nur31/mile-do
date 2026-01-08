package domain

import (
	"time"

	repo "github.com/ali-nur31/mile-do/internal/db"
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
