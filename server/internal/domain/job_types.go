package domain

import (
	"encoding/json"
	"log/slog"

	"github.com/hibiken/asynq"
)

const (
	TypeGenerateDefaultGoals             = "generate:default:goals"
	TypeGenerateRecurringTasks           = "generate:recurring:tasks"
	TypeGenerateRecurringTasksByTemplate = "generate:recurring:tasks:by:template"
	TypeDeleteRecurringTasksByTemplateID = "delete:recurring:tasks:by:template:id"
)

func NewGenerateDefaultGoals(id int32) *asynq.Task {
	encodedPayload, err := json.Marshal(id)
	if err != nil {
		slog.Error("couldn't convert map to bytes", "error", err)
		return nil
	}

	return asynq.NewTask(TypeGenerateDefaultGoals, encodedPayload)
}

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
	encodedPayload, err := json.Marshal(id)
	if err != nil {
		slog.Error("couldn't convert map to bytes", "error", err)
		return nil
	}

	return asynq.NewTask(TypeDeleteRecurringTasksByTemplateID, encodedPayload)
}
