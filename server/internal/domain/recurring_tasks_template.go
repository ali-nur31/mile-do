package domain

import (
	"context"
	"time"

	"github.com/ali-nur31/mile-do/internal/repository/db"
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
	if len(templates) == 0 {
		return nil
	}

	output := make([]RecurringTasksTemplateOutput, len(templates))
	for i, t := range templates {
		output[i] = *ToRecurringTasksTemplateOutput(&t)
	}
	return output
}

type RecurringTasksTemplateService interface {
	ListRecurringTasksTemplates(ctx context.Context, userId int32) ([]RecurringTasksTemplateOutput, error)
	GetRecurringTasksTemplateByID(ctx context.Context, id int64, userId int32) (*RecurringTasksTemplateOutput, error)
	CreateRecurringTasksTemplate(ctx context.Context, input CreateRecurringTasksTemplateInput) (*RecurringTasksTemplateOutput, error)
	UpdateRecurringTasksTemplateByID(ctx context.Context, dbTemplate RecurringTasksTemplateOutput, updatingTemplate UpdateRecurringTasksTemplateInput) (*RecurringTasksTemplateOutput, error)
	DeleteRecurringTasksTemplateByID(ctx context.Context, id int64, userId int32) error
	ListRecurringTasksTemplatesDueForGeneration(ctx context.Context, qtx repo.Querier) ([]RecurringTasksTemplateOutput, error)
	UpdateLastGeneratedDateInRecurringTasksTemplateByID(ctx context.Context, qtx repo.Querier, updatingTemplate UpdateLastGeneratedDateInRecurringTasksTemplateInput) error
}
