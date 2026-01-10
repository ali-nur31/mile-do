package workers

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/ali-nur31/mile-do/internal/domain"
	"github.com/ali-nur31/mile-do/internal/service"
	"github.com/hibiken/asynq"
)

type RecurringTasksTemplatesWorker struct {
	service service.TaskService
}

func NewRecurringTasksTemplatesWorker(service service.TaskService) *RecurringTasksTemplatesWorker {
	return &RecurringTasksTemplatesWorker{service: service}
}

func (w *RecurringTasksTemplatesWorker) GenerateRecurringTasks(ctx context.Context, t *asynq.Task) error {
	slog.Info("executing recurring tasks generation job")

	err := w.service.CreateTasksByRecurringTasksTemplates(ctx)
	if err != nil {
		slog.Error("failed to execute recurring tasks generation job", "error", err)
		return err
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

	err = w.service.CreateTasksByRecurringTasksTemplate(ctx, template)
	if err != nil {
		slog.Error("failed to execute recurring tasks generation by template job", "error", err)
		return err
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
