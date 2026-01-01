package auth

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/ali-nur31/mile-do/config"
	"github.com/golang-jwt/jwt/v5"
)

type JwtManager struct {
	jwt *config.Jwt
}

type AccessClaims struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
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

func (m *JwtManager) CreateToken(id int64, email string) (TokensData, error) {
	accessClaims := AccessClaims{
		ID:    id,
		Email: email,
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

func (m *JwtManager) VerifyToken(tokenString string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.jwt.AccessKey), nil
	})

	if err != nil {
		slog.Error("failed to parse token", "error", err)
		return nil, err
	} else if claims, ok := token.Claims.(*AccessClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token claims")
	}
}
