package domain

import (
	"encoding/json"
	"log/slog"

	"github.com/hibiken/asynq"
)

const (
	TypeGenerateRecurringTasks           = "generate:recurring:tasks"
	TypeGenerateRecurringTasksByTemplate = "generate:recurring:tasks:by:template"
	TypeDeleteRecurringTasksByTemplateID = "delete:recurring:tasks:by:template:id"
)

func NewGenerateRecurringTasksTask() *asynq.Task {
	return asynq.NewTask(TypeGenerateRecurringTasks, []byte{})
}

func NewGenerateRecurringTasksByTemplateTask(template *RecurringTasksTemplateOutput) *asynq.Task {
	encodedPayload, err := json.Marshal(template)
	if err != nil {
		slog.Error("couldn't convert map to bytes", "error", err)
		return nil
	}

	return asynq.NewTask(TypeGenerateRecurringTasksByTemplate, encodedPayload)
}

func NewDeleteRecurringTasksByTemplateIDTask(id int64) *asynq.Task {
	payload := map[string]interface{}{"id": id}

	encodedPayload, err := json.Marshal(payload)
	if err != nil {
		slog.Error("couldn't convert map to bytes", "error", err)
		return nil
	}

	return asynq.NewTask(TypeDeleteRecurringTasksByTemplateID, encodedPayload)
}
