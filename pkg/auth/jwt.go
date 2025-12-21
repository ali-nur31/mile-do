package auth

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtManager struct {
	secretKey string
}

type AuthClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func NewJwtManager(secretKey string) (*JwtManager, error) {
	if secretKey == "" {
		slog.Error("empty secretKey for jwt")
		return nil, errors.New("empty secretKey for jwt")
	}
	return &JwtManager{
		secretKey: secretKey,
	}, nil
}

func (m *JwtManager) CreateToken(email string) (string, error) {
	claims := AuthClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenString, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (m *JwtManager) VerifyToken(tokenString string) (*AuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.secretKey), nil
	})

	if err != nil {
		slog.Error("failed to parse token", "error", err)
		return nil, err
	} else if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token claims")
	}
}
