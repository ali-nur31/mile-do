package middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ali-nur31/mile-do/internal/service"
	"github.com/ali-nur31/mile-do/pkg/auth"
	"github.com/labstack/echo/v4"
)

type AuthTokenManager interface {
	VerifyToken(tokenString, tokenType string) (*auth.Claims, error)
}

type AuthMiddleware struct {
	tokenManager        AuthTokenManager
	refreshTokenService service.RefreshTokenService
}

func NewAuthMiddleware(tokenManager AuthTokenManager, refreshTokenService service.RefreshTokenService) *AuthMiddleware {
	return &AuthMiddleware{
		tokenManager:        tokenManager,
		refreshTokenService: refreshTokenService,
	}
}

func (m *AuthMiddleware) TokenCheckMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing authorization header"})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid authorization format"})
			}

			tokenString := parts[1]

			claims, err := m.tokenManager.VerifyToken(tokenString, "access")
			if err != nil {
				slog.Error("couldn't verify token", "error", err)
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid token", "error": err.Error()})
			} else if errors.Is(err, auth.TokenExpiredError) {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
			}

			c.Set("userId", claims.ID)

			return next(c)
		}
	}
}
