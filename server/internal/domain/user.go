package domain

import (
	"time"

	"github.com/ali-nur31/mile-do/internal/repository/db"
)

type UserOutput struct {
	ID           int64
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

func ToUserOutput(u *repo.User) *UserOutput {
	return &UserOutput{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt.Time,
	}
}
