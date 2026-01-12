package domain

import (
	"context"
	"time"

	"github.com/ali-nur31/mile-do/internal/repository/db"
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

type AuthTokenManager interface {
	CreateTokens(id int64) (*auth.TokensData, error)
	VerifyToken(tokenString, tokenType string) (*auth.Claims, error)
}

type AuthPasswordManager interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type AuthService interface {
	RegisterUser(ctx context.Context, user UserInput) (*AuthOutput, error)
	LoginUser(ctx context.Context, user UserInput) (*AuthOutput, error)
	LogoutUser(ctx context.Context, userId int32) error
	RefreshTokens(ctx context.Context, refreshToken string) (*AuthOutput, error)
}

type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (*UserOutput, error)
	GetUserByID(ctx context.Context, id int64) (*UserOutput, error)
	CreateUser(ctx context.Context, qtx repo.Querier, user UserInput) (*UserOutput, error)
}
