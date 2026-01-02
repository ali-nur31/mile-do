package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/ali-nur31/mile-do/config"
	_ "github.com/ali-nur31/mile-do/docs"
	repo "github.com/ali-nur31/mile-do/internal/db"
	"github.com/ali-nur31/mile-do/internal/service"
	"github.com/ali-nur31/mile-do/internal/transport/http/middleware"
	v1 "github.com/ali-nur31/mile-do/internal/transport/http/v1"
	"github.com/ali-nur31/mile-do/internal/worker"
	"github.com/ali-nur31/mile-do/internal/worker/jobs"
	"github.com/ali-nur31/mile-do/pkg/asynq_jobs"
	"github.com/ali-nur31/mile-do/pkg/auth"
	"github.com/ali-nur31/mile-do/pkg/logger"
	"github.com/ali-nur31/mile-do/pkg/postgres"
	"github.com/ali-nur31/mile-do/pkg/redis_db"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

// @title           Mile-Do API
// @version         1.0
// @description     API for Mile-Do, simple clone of TickTick.
// @termsOfService  http://swagger.io/terms/

// @contact.name    API Support
// @contact.email   support@swagger.io

// @host            localhost:8080
// @BasePath        /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token
func main() {
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
	}

	logger.InitializeLogger()

	cfg := config.MustLoad()

	pg, err := postgres.InitializeDatabaseConnection(ctx, &cfg.DB)
	if err != nil {
		os.Exit(1)
	}

	defer pg.Pool.Close(ctx)

	queries := repo.New(pg.Pool)

	_, err = redis_db.InitializeRedisConnection(ctx, &cfg.Redis)
	if err != nil {
		os.Exit(1)
	}

	asynq, err := asynq_jobs.InitializeAsynqClient(&cfg.Redis)
	if err != nil {
		os.Exit(1)
	}

	defer asynq.Client.Close()

	passwordManager := auth.NewBcryptPasswordManager()

	jwtTokenManager, err := auth.NewJwtManager(&cfg.Jwt)
	if err != nil {
		os.Exit(1)
	}

	refreshTokenService := service.NewRefreshTokenService(queries)

	authMiddleware := middleware.NewAuthMiddleware(jwtTokenManager, refreshTokenService)

	userService := service.NewUserService(queries, jwtTokenManager, refreshTokenService, passwordManager)
	authHandler := v1.NewAuthHandler(userService)
	userHandler := v1.NewUserHandler(userService)

	goalService := service.NewGoalService(queries)
	goalHandler := v1.NewGoalHandler(goalService)

	taskService := service.NewTaskService(queries)
	taskHandler := v1.NewTaskHandler(taskService)

	router := v1.NewRouter(
		*authMiddleware,
		*authHandler,
		*userHandler,
		*goalHandler,
		*taskHandler,
	)

	e := echo.New()

	apiGroup := e.Group("api/v1")

	router.InitRoutes(apiGroup)

	taskGenerateRecurringJob := jobs.NewTaskGenerateRecurringJob(taskService)

	backgroundWorker := worker.NewWorker(&cfg.Redis, taskGenerateRecurringJob)

	go func() {
		if err = backgroundWorker.Run(); err != nil {
			slog.Error("failed to run worker, exit", "error", err)
			os.Exit(1)
		}
	}()

	port := cfg.Api.Port
	if err := e.Start(port); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}
