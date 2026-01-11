package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	repo "github.com/ali-nur31/mile-do/internal/db"
	"github.com/ali-nur31/mile-do/internal/domain"
	"github.com/ali-nur31/mile-do/internal/service"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RecurringTasksTemplatesWorker struct {
	pool    *pgxpool.Pool
	service service.TaskService
}

func NewRecurringTasksTemplatesWorker(pool *pgxpool.Pool, service service.TaskService) *RecurringTasksTemplatesWorker {
	return &RecurringTasksTemplatesWorker{
		pool:    pool,
		service: service,
	}
}

func (w *RecurringTasksTemplatesWorker) GenerateRecurringTasksDueForGeneration(ctx context.Context, t *asynq.Task) error {
	slog.Info("executing recurring tasks generation job")

	tx, err := w.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(context.Background())
	}()

	qtx := repo.New(tx)

	err = w.service.CreateTasksByRecurringTasksTemplatesDueForGeneration(ctx, qtx)
	if err != nil {
		slog.Error("failed to execute recurring tasks generation job", "error", err)
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("couldn't commit transaction for generate recurring tasks: %w", err)
	}

	slog.Info("ended execution of recurring tasks generation job")
	return nil
}

func (w *RecurringTasksTemplatesWorker) GenerateRecurringTasksByTemplate(ctx context.Context, t *asynq.Task) error {
	slog.Info("executing recurring tasks generation by template job")

	templateBytes := t.Payload()
	var template domain.RecurringTasksTemplateOutput
	err := json.Unmarshal(templateBytes, &template)
	if err != nil {
		slog.Error("Error unmarshalling bytes to map", "error", err)
		return err
	}

	tx, err := w.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(context.Background())
	}()

	qtx := repo.New(tx)

	err = w.service.CreateTasksByRecurringTasksTemplate(ctx, qtx, template)
	if err != nil {
		slog.Error("failed to execute recurring tasks generation by template job", "error", err)
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("couldn't commit transaction for generate recurring tasks by template: %w", err)
	}

	slog.Info("ended execution of recurring tasks generation by template job")
	return nil
}

func (w *RecurringTasksTemplatesWorker) DeleteRecurringTasksByTemplateID(ctx context.Context, t *asynq.Task) error {
	slog.Info("executing recurring tasks deletion by template id job")

	templateIdBytes := t.Payload()
	var templateId int64
	err := json.Unmarshal(templateIdBytes, &templateId)
	if err != nil {
		slog.Error("Error unmarshalling bytes to map", "error", err)
		return err
	}

	err = w.service.DeleteFutureTasksByRecurringTasksTemplateID(ctx, templateId)
	if err != nil {
		slog.Error("failed to execute recurring tasks generation job", "error", err)
		return err
	}

	slog.Info("ended execution of recurring tasks deletion by template id job")
	return nil
}
