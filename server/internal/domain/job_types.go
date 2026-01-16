package domain

import (
	"encoding/json"
	"log/slog"

	"github.com/hibiken/asynq"
)

const (
	TypeGenerateDefaultGoals                   = "generate:default:goals"
	TypeGenerateRecurringTasksDueForGeneration = "generate:recurring:tasks:due:for:generation"
	TypeGenerateRecurringTasksByTemplate       = "generate:recurring:tasks:by:template"
	TypeDeleteRecurringTasksByTemplateID       = "delete:recurring:tasks:by:template:id"
)

func NewGenerateDefaultGoalsTask(id int32) *asynq.Task {
	encodedPayload, err := json.Marshal(id)
	if err != nil {
		slog.Error("couldn't convert map to bytes", "error", err)
		return nil
	}

	return asynq.NewTask(TypeGenerateDefaultGoals, encodedPayload)
}

func NewGenerateRecurringTasksDueForGenerationTask() *asynq.Task {
	return asynq.NewTask(TypeGenerateRecurringTasksDueForGeneration, []byte{})
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
