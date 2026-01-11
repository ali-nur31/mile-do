package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ali-nur31/mile-do/config"
	_ "github.com/ali-nur31/mile-do/docs"
	repo "github.com/ali-nur31/mile-do/internal/db"
	"github.com/ali-nur31/mile-do/internal/jobs"
	"github.com/ali-nur31/mile-do/internal/jobs/workers"
	"github.com/ali-nur31/mile-do/internal/service"
	"github.com/ali-nur31/mile-do/internal/transport/http/middleware"
	v1 "github.com/ali-nur31/mile-do/internal/transport/http/v1"
	"github.com/ali-nur31/mile-do/pkg/asynq_jobs"
	"github.com/ali-nur31/mile-do/pkg/auth"
	"github.com/ali-nur31/mile-do/pkg/logger"
	"github.com/ali-nur31/mile-do/pkg/postgres"
	"github.com/ali-nur31/mile-do/pkg/redis_db"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron/v3"
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
		slog.Error("couldn't initialize database pool connection", "error", err)
		os.Exit(1)
	}

	defer pg.Pool.Close()

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

	userService := service.NewUserService(queries, passwordManager)
	userHandler := v1.NewUserHandler(userService)

	authService := service.NewAuthService(queries, asynq.Client, pg.Pool, userService, jwtTokenManager, refreshTokenService, passwordManager)
	authHandler := v1.NewAuthHandler(authService)

	goalService := service.NewGoalService(queries)
	goalHandler := v1.NewGoalHandler(goalService)

	recurringTasksTemplateService := service.NewRecurringTasksTemplateService(queries, asynq.Client)
	recurringTasksTemplateHandler := v1.NewRecurringTasksTemplateHandler(recurringTasksTemplateService)

	taskService := service.NewTaskService(queries, pg.Pool, recurringTasksTemplateService)
	taskHandler := v1.NewTaskHandler(taskService)

	router := v1.NewRouter(
		cfg.Redis,
		*authMiddleware,
		*authHandler,
		*userHandler,
		*goalHandler,
		*recurringTasksTemplateHandler,
		*taskHandler,
	)

	e := echo.New()

	e.Use(echo_middleware.CORSWithConfig(echo_middleware.CORSConfig{
		AllowOrigins:     []string{cfg.Api.FrontendUrl},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PATCH, echo.DELETE, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAuthorization, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	apiGroup := e.Group("api/v1")

	router.InitRoutes(apiGroup)

	c := cron.New()

	scheduler := service.NewScheduler(c, asynq.Client)

	scheduler.InitSchedules()

	c.Start()
	slog.Info("Cron scheduler started")

	goalsWorker := workers.NewGoalsWorker(pg.Pool, goalService)
	recurringTasksTemplatesWorker := workers.NewRecurringTasksTemplatesWorker(pg.Pool, taskService)

	backgroundWorker := jobs.NewJobRouter(&cfg.Redis, goalsWorker, recurringTasksTemplatesWorker)

	go func() {
		if err = backgroundWorker.Run(); err != nil {
			slog.Error("failed to run jobs, exit", "error", err)
			os.Exit(1)
		}
	}()

	go func() {
		port := cfg.Api.Port
		if err = e.Start(port); err != nil {
			slog.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	slog.Info("Received shutdown signal. shutting down...")

	ctxCron := c.Stop()
	<-ctxCron.Done()
	slog.Info("Cron scheduler stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	slog.Info("Server exited properly")
}
