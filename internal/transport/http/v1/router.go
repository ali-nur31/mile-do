package v1

import (
	"github.com/ali-nur31/mile-do/internal/transport/http/middleware"
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
)

type Router struct {
	authMiddleware middleware.AuthMiddleware
	authHandler    AuthHandler
	userHandler    UserHandler
	goalHandler    GoalHandler
}

func NewRouter(
	authMiddleware middleware.AuthMiddleware,
	authHandler AuthHandler,
	userHandler UserHandler,
	goalHandler GoalHandler,
) *Router {
	return &Router{
		authMiddleware: authMiddleware,
		authHandler:    authHandler,
		userHandler:    userHandler,
		goalHandler:    goalHandler,
	}
}

func (r Router) InitRoutes(api *echo.Group) {
	api.GET("/swagger/*", echoSwagger.WrapHandler)

	auth := api.Group("/auth")
	{
		auth.POST("/register", r.authHandler.RegisterUser)
		auth.POST("/login", r.authHandler.LoginUser)
	}

	users := api.Group("/users")
	users.Use(r.authMiddleware.TokenCheckMiddleware())
	{
		users.GET("/me", r.userHandler.GetUserByEmail)
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
