package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/ali-nur31/mile-do/config"
	repo "github.com/ali-nur31/mile-do/internal/db"
	"github.com/ali-nur31/mile-do/internal/service"
	"github.com/ali-nur31/mile-do/internal/transport/http/middleware"
	v1 "github.com/ali-nur31/mile-do/internal/transport/http/v1"
	"github.com/ali-nur31/mile-do/pkg/auth"
	"github.com/ali-nur31/mile-do/pkg/logger"
	"github.com/ali-nur31/mile-do/pkg/postgres"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

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
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	defer pg.Pool.Close(ctx)

	queries := repo.New(pg.Pool)

	passwordManager := auth.NewBcryptPasswordManager()

	jwtTokenManager, err := auth.NewJwtManager(cfg.Api.JWTSecretKey)
	if err != nil {
		os.Exit(1)
	}

	authMiddleware := middleware.NewAuthMiddleware(jwtTokenManager)

	userService := service.NewUserService(queries, jwtTokenManager, passwordManager)
	authHandler := v1.NewAuthHandler(userService)
	userHandler := v1.NewUserHandler(userService)

	goalService := service.NewGoalService(queries)
	goalHandler := v1.NewGoalHandler(goalService)

	router := v1.NewRouter(
		*authMiddleware,
		*authHandler,
		*userHandler,
		*goalHandler,
	)

	e := echo.New()

	apiGroup := e.Group("api/v1")

	router.InitRoutes(apiGroup)

	port := cfg.Api.Port
	if err := e.Start(port); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}
