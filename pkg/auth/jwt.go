package auth

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/ali-nur31/mile-do/config"
	"github.com/golang-jwt/jwt/v5"
)

var TokenExpiredError = errors.New("token has expired")

type JwtManager struct {
	jwt *config.Jwt
}

type AccessClaims struct {
	ID int64 `json:"id"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	ID int64 `json:"id"`
	jwt.RegisteredClaims
}

type TokensData struct {
	AccessToken     string
	AccessTokenExp  time.Time
	RefreshToken    string
	RefreshTokenExp time.Time
}

func NewJwtManager(jwt *config.Jwt) (*JwtManager, error) {
	return &JwtManager{
		jwt: jwt,
	}, nil
}

func (m *JwtManager) CreateTokens(id int64) (TokensData, error) {
	accessClaims := AccessClaims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(m.jwt.AccessExpMins))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshClaims := RefreshClaims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * time.Duration(m.jwt.RefreshExpDays))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		accessClaims,
	)

	refreshToken := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		refreshClaims,
	)

	accessTokenString, err := accessToken.SignedString([]byte(m.jwt.AccessKey))
	if err != nil {
		return TokensData{}, err
	}

	refreshTokenString, err := refreshToken.SignedString([]byte(m.jwt.RefreshKey))
	if err != nil {
		return TokensData{}, err
	}

	return TokensData{
		AccessToken:     accessTokenString,
		AccessTokenExp:  accessClaims.ExpiresAt.Time,
		RefreshToken:    refreshTokenString,
		RefreshTokenExp: refreshClaims.ExpiresAt.Time,
	}, nil
}

func (m *JwtManager) VerifyAccessToken(tokenString string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.jwt.AccessKey), nil
	})

	if err != nil {
		slog.Error("failed to parse access token", "error", err)
		return nil, err
	} else if claims, ok := token.Claims.(*AccessClaims); ok && token.Valid {
		return claims, nil
	} else if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, TokenExpiredError
	}

	return nil, fmt.Errorf("invalid access token claims")
}

func (m *JwtManager) VerifyRefreshToken(tokenString string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.jwt.RefreshKey), nil
	})

	if err != nil {
		slog.Error("failed to parse refresh token", "error", err)
		return nil, err
	} else if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		return claims, nil
	} else if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, TokenExpiredError
	}

	return nil, fmt.Errorf("invalid refresh token claims")
}
