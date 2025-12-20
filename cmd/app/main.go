package main

import (
	"context"
	"log/slog"

	"github.com/ali-nur31/mile-do/config"
	repo "github.com/ali-nur31/mile-do/internal/db"
	"github.com/ali-nur31/mile-do/internal/service"
	v1 "github.com/ali-nur31/mile-do/internal/transport/http/v1"
	"github.com/ali-nur31/mile-do/pkg/logger"
	"github.com/ali-nur31/mile-do/pkg/postgres"
	"github.com/labstack/echo/v4"
)

func main() {
	ctx := context.Background()

	logger.InitializeLogger()

	cfg := config.MustLoad()

	pool := postgres.InitializeDatabaseConnection(ctx, &cfg.DB)

	queries := repo.New(pool)

	userService := service.NewUserService(queries)
	userHandler := v1.NewUserHandler(userService)

	router := v1.NewRouter(
		*userHandler,
	)

	e := echo.New()

	apiGroup := e.Group("api/v1")

	router.InitRoutes(apiGroup)

	port := cfg.Api.Port
	if err := e.Start(port); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}
