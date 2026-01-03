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

func NewGenerateRecurringTasksByTemplateTask(template RecurringTasksTemplateOutput) *asynq.Task {
	payload := map[string]interface{}{
		"id":                  template.ID,
		"user_id":             template.UserID,
		"goal_id":             template.GoalID,
		"title":               template.Title,
		"scheduled_datetime":  template.ScheduledDatetime,
		"has_time":            template.HasTime,
		"duration_minutes":    template.DurationMinutes,
		"recurrence_rrule":    template.RecurrenceRrule,
		"last_generated_date": template.LastGeneratedDate,
		"created_at":          template.CreatedAt,
	}

	encodedPayload, err := json.Marshal(payload)
	if err != nil {
		slog.Error("couldn't convert map to bytes", "error", err)
		return nil
	}

	return asynq.NewTask(TypeDeleteRecurringTasksByTemplateID, encodedPayload)
}

func NewDeleteRecurringTasksByTemplateIDTask(id int64) *asynq.Task {
	payload := map[string]interface{}{"id": id}

	encodedPayload, err := json.Marshal(payload)
	if err != nil {
		slog.Error("couldn't convert map to bytes", "error", err)
		return nil
	}

	return asynq.NewTask(TypeGenerateRecurringTasks, encodedPayload)
}
