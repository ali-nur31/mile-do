package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/ali-nur31/mile-do/config"
	"github.com/golang-jwt/jwt/v5"
)

var TokenExpiredError = errors.New("token has expired")

type JwtManager struct {
	jwt *config.Jwt
}

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

func NewJwtManager(jwt *config.Jwt) (*JwtManager, error) {
	return &JwtManager{
		jwt: jwt,
	}, nil
}

func (m *JwtManager) CreateTokens(id int64) (*TokensData, error) {
	accessClaims := Claims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(m.jwt.AccessExpMins))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshClaims := Claims{
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
		return nil, fmt.Errorf("couldn't sign access token: %w", err)
	}

	refreshTokenString, err := refreshToken.SignedString([]byte(m.jwt.RefreshKey))
	if err != nil {
		return nil, fmt.Errorf("couldn't sign refresh token: %w", err)
	}

	return &TokensData{
		AccessToken:     accessTokenString,
		AccessTokenExp:  accessClaims.ExpiresAt.Time,
		RefreshToken:    refreshTokenString,
		RefreshTokenExp: refreshClaims.ExpiresAt.Time,
	}, nil
}

func (m *JwtManager) VerifyToken(tokenString, tokenType string) (*Claims, error) {
	var secretKey string
	if tokenType == "access" {
		secretKey = m.jwt.AccessKey
	} else if tokenType == "refresh" {
		secretKey = m.jwt.RefreshKey
	} else {
		return nil, fmt.Errorf("invalid tokenType param: %v", tokenType)
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("couldn't parse %v token: %w", tokenType, err)
	} else if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, TokenExpiredError
	}

	return nil, fmt.Errorf("invalid %v token claims", tokenType)
}
