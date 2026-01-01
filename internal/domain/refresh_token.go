package domain

import "time"

type CreateRefreshTokenInput struct {
	UserID    int32
	TokenHash string
	ExpiresAt time.Time
}
