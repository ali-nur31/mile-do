package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/ali-nur31/taskee/config"
	"github.com/jackc/pgx/v5"
)

func InitializeDatabaseConnection(ctx context.Context) {
	cfg, err := config.MustLoad()
	if err != nil {
		slog.Error("couldn't get environment variables", "error", err)
		os.Exit(1)
	}

	conn, err := pgx.Connect(ctx,
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s",
			cfg.Username,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Name,
		),
	)
	if err != nil {
		slog.Error("couldn't connect to database", "error", err)
		os.Exit(1)
	}

	fmt.Println("Connection to database is established")

	defer conn.Close(ctx)
}
