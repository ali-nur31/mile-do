package domain

import (
	"time"

	repo "github.com/ali-nur31/mile-do/internal/db"
)

type CreateRecurringTasksTemplateInput struct {
	UserID            int32
	GoalID            int32
	Title             string
	ScheduledDatetime time.Time
	HasTime           bool
	DurationMinutes   int32
	RecurrenceRrule   string
}

type UpdateRecurringTasksTemplateInput struct {
	ID                int64
	UserID            int32
	GoalID            int32
	Title             string
	ScheduledDatetime time.Time
	HasTime           bool
	DurationMinutes   int32
	RecurrenceRrule   string
}

type UpdateLastGeneratedDateInRecurringTasksTemplateInput struct {
	ID                int64
	LastGeneratedDate time.Time
}

type RecurringTasksTemplateOutput struct {
	ID                int64
	UserID            int32
	GoalID            int32
	Title             string
	ScheduledDatetime time.Time
	HasTime           bool
	DurationMinutes   int32
	RecurrenceRrule   string
	LastGeneratedDate time.Time
	CreatedAt         time.Time
}

func ToRecurringTasksTemplateOutput(template *repo.RecurringTasksTemplate) *RecurringTasksTemplateOutput {
	return &RecurringTasksTemplateOutput{
		ID:                template.ID,
		UserID:            template.UserID,
		GoalID:            template.GoalID,
		Title:             template.Title,
		ScheduledDatetime: template.ScheduledDatetime.Time,
		HasTime:           template.HasTime,
		DurationMinutes:   template.DurationMinutes,
		RecurrenceRrule:   template.RecurrenceRrule,
		LastGeneratedDate: template.LastGeneratedDate.Time,
		CreatedAt:         template.CreatedAt.Time,
	}
}

func ToRecurringTasksTemplateOutputList(templates []repo.RecurringTasksTemplate) []RecurringTasksTemplateOutput {
	output := make([]RecurringTasksTemplateOutput, len(templates))
	for i, t := range templates {
		output[i] = *ToRecurringTasksTemplateOutput(&t)
	}
	return output
}
