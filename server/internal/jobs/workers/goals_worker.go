package workers

import (
	"context"
	"encoding/json"
	"log/slog"

	repo "github.com/ali-nur31/mile-do/internal/db"
	"github.com/ali-nur31/mile-do/internal/domain"
	"github.com/ali-nur31/mile-do/internal/service"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GoalsWorker struct {
	service service.GoalService
	pool    *pgxpool.Pool
}

func NewGoalsWorker(service service.GoalService, pool *pgxpool.Pool) *GoalsWorker {
	return &GoalsWorker{
		service: service,
		pool:    pool,
	}
}

func (w *GoalsWorker) GenerateDefaultTasks(ctx context.Context, t *asynq.Task) error {
	slog.Info("executing default tasks generation job")

	userIdBytes := t.Payload()
	var userId int32
	err := json.Unmarshal(userIdBytes, &userId)
	if err != nil {
		slog.Error("Error unmarshalling bytes to map", "error", err)
		return err
	}

	defaultGoals := []domain.CreateGoalInput{
		{
			UserID:       userId,
			Title:        "Routine",
			Color:        "#73260A",
			CategoryType: "maintenance",
		},
		{
			UserID:       userId,
			Title:        "Other",
			Color:        "#0096ff",
			CategoryType: "other",
		},
	}

	tx, err := w.pool.Begin(ctx)
	if err != nil {
		slog.Error("failed to execute default tasks generation job, transaction begin error", "error", err)
		return err
	}
	defer func() {
		_ = tx.Rollback(context.Background())
	}()

	qtx := repo.New(tx)

	for _, input := range defaultGoals {
		_, err = w.service.CreateGoal(ctx, qtx, input)
		if err != nil {
			slog.Error("failed to execute default tasks generation job", "error", err)
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		slog.Error("couldn't commit transaction for default tasks generation job", "error", err)
	}

	return nil
}
