package jobs

import (
	"context"
	"log/slog"

	"github.com/ali-nur31/mile-do/internal/service"
	"github.com/hibiken/asynq"
)

type TaskGenerateRecurringJob struct {
	service service.TaskService
}

func NewTaskGenerateRecurringJob(service service.TaskService) *TaskGenerateRecurringJob {
	return &TaskGenerateRecurringJob{service: service}
}

func (j *TaskGenerateRecurringJob) Process(ctx context.Context, t *asynq.Task) error {
	slog.Info("executing recurring task generation job")
	return nil
}
