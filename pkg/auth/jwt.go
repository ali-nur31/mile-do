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
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		},
	)

	tokenString, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (m *JwtManager) VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.secretKey), nil
	})

	if err != nil {
		slog.Error("failed to verify token", "error", err)
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
