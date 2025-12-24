package middleware

import (
	"net/http"
	"strings"

	"github.com/ali-nur31/mile-do/pkg/auth"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	tokenManager auth.JwtManager
}

func NewAuthMiddleware(tokenManager auth.JwtManager) *AuthMiddleware {
	return &AuthMiddleware{
		tokenManager: tokenManager,
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

			claims, err := m.tokenManager.VerifyToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid or expired token"})
			}

			c.Set("userId", claims.ID)
			c.Set("email", claims.Email)

			return next(c)
		}
	}
}
