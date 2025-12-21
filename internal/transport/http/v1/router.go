package v1

import (
	"github.com/ali-nur31/mile-do/internal/transport/http/middleware"
	"github.com/labstack/echo/v4"
)

type Router struct {
	authHandler    AuthHandler
	authMiddleware middleware.AuthMiddleware
}

func NewRouter(
	authHandler AuthHandler,
	authMiddleware middleware.AuthMiddleware,
) *Router {
	return &Router{
		authHandler:    authHandler,
		authMiddleware: authMiddleware,
	}
}

func (r Router) InitRoutes(api *echo.Group) {
	authPublic := api.Group("/auth")
	{
		authPublic.POST("/register", r.authHandler.RegisterUser)
		authPublic.POST("/login", r.authHandler.LoginUser)
	}

	authPrivate := api.Group("/auth")
	authPrivate.Use(r.authMiddleware.TokenCheckMiddleware())
	{
		authPrivate.GET("/me", r.authHandler.GetUserByEmail)
	}
}
