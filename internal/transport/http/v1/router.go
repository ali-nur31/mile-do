package v1

import (
	"github.com/ali-nur31/mile-do/internal/transport/http/middleware"
	"github.com/labstack/echo/v4"
)

type Router struct {
	authMiddleware middleware.AuthMiddleware
	authHandler    AuthHandler
	goalHandler    GoalHandler
}

func NewRouter(
	authMiddleware middleware.AuthMiddleware,
	authHandler AuthHandler,
	goalHandler GoalHandler,
) *Router {
	return &Router{
		authMiddleware: authMiddleware,
		authHandler:    authHandler,
		goalHandler:    goalHandler,
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

	goals := api.Group("/goals")
	goals.Use(r.authMiddleware.TokenCheckMiddleware())
	{
		goals.GET("/", r.goalHandler.GetGoals)
		goals.GET("/:id", r.goalHandler.GetGoalByID)
		goals.POST("/", r.goalHandler.CreateGoal)
		goals.PUT("/", r.goalHandler.UpdateGoal)
		goals.DELETE("/:id", r.goalHandler.DeleteGoalByID)
	}
}
