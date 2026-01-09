package domain

import (
	"time"

	repo "github.com/ali-nur31/mile-do/internal/db"
	"github.com/ali-nur31/mile-do/pkg/auth"
)

type UserInput struct {
	Email    string
	Password string
}

type AuthOutput struct {
	AccessToken  string
	RefreshToken string
}

func ToAuthOutput(t *auth.TokensData) *AuthOutput {
	return &AuthOutput{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
	}
}

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
