package middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ali-nur31/mile-do/internal/domain"
	"github.com/ali-nur31/mile-do/pkg/auth"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	authCacheRepo       domain.AuthCacheRepo
	tokenManager        domain.AuthTokenManager
	refreshTokenService domain.RefreshTokenService
}

func NewAuthMiddleware(authCacheRepo domain.AuthCacheRepo, tokenManager domain.AuthTokenManager, refreshTokenService domain.RefreshTokenService) *AuthMiddleware {
	return &AuthMiddleware{
		authCacheRepo:       authCacheRepo,
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

			isBlocked, err := m.authCacheRepo.IsTokenBlocked(c.Request().Context(), tokenString)
			if err != nil {
				slog.Error("couldn't check if token is blocked", "error", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
			}
			if isBlocked {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "unauthorized"})
			}

			claims, err := m.tokenManager.VerifyToken(tokenString, "access")
			if err != nil {
				slog.Error("couldn't verify token", "error", err)
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid token", "error": err.Error()})
			} else if errors.Is(err, auth.TokenExpiredError) {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
			}

			c.Set("claims", claims)
			c.Set("accessToken", tokenString)

			return next(c)
		}
	}
}
