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
	taskHandler    TaskHandler
}

func NewRouter(
	authMiddleware middleware.AuthMiddleware,
	authHandler AuthHandler,
	userHandler UserHandler,
	goalHandler GoalHandler,
	taskHandler TaskHandler,
) *Router {
	return &Router{
		authMiddleware: authMiddleware,
		authHandler:    authHandler,
		userHandler:    userHandler,
		goalHandler:    goalHandler,
		taskHandler:    taskHandler,
	}
}

func (r Router) InitRoutes(api *echo.Group) {
	api.GET("/swagger/*", echoSwagger.WrapHandler)

	auth := api.Group("/auth")
	{
		auth.POST("/register", r.authHandler.RegisterUser)
		auth.POST("/login", r.authHandler.LoginUser)
		auth.POST("/refresh", r.authHandler.RefreshAccessToken)
	}

	users := api.Group("/users")
	users.Use(r.authMiddleware.TokenCheckMiddleware())
	{
		users.GET("/me", r.userHandler.GetUserByEmail)
		users.DELETE("/", r.userHandler.LogoutUser)
	}

	goals := api.Group("/goals")
	goals.Use(r.authMiddleware.TokenCheckMiddleware())
	{
		goals.GET("/", r.goalHandler.GetGoals)
		goals.GET("/:id", r.goalHandler.GetGoalByID)
		goals.GET("/:id/tasks", r.taskHandler.GetTasksByGoalID)
		goals.POST("/", r.goalHandler.CreateGoal)
		goals.PATCH("/", r.goalHandler.UpdateGoal)
		goals.DELETE("/:id", r.goalHandler.DeleteGoalByID)
	}

	tasks := api.Group("/tasks")
	tasks.Use(r.authMiddleware.TokenCheckMiddleware())
	{
		tasks.GET("/", r.taskHandler.GetTasks)
		tasks.GET("/inbox", r.taskHandler.GetInboxTasks)
		tasks.GET("/period", r.taskHandler.GetTasksByPeriod)
		tasks.GET("/analyze", r.taskHandler.AnalyzeForToday)
		tasks.GET("/:id", r.taskHandler.GetTaskByID)
		tasks.POST("/", r.taskHandler.CreateTask)
		tasks.PATCH("/:id", r.taskHandler.UpdateTask)
		tasks.DELETE("/:id", r.taskHandler.DeleteTaskByID)
	}
}
