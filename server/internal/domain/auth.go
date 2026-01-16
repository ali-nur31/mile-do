package domain

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var TokenExpiredError = errors.New("token has expired")

type Claims struct {
	ID int64 `json:"id"`
	jwt.RegisteredClaims
}

type TokensData struct {
	AccessToken     string
	AccessTokenExp  time.Time
	RefreshToken    string
	RefreshTokenExp time.Time
}

type AuthInput struct {
	Email    string
	Password string
}

type AuthOutput struct {
	AccessToken  string
	RefreshToken string
}

func ToAuthOutput(t *TokensData) *AuthOutput {
	return &AuthOutput{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
	}
}
