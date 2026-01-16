package dto

import (
	"github.com/ali-nur31/mile-do/internal/domain"
)

type UpdateRecurringTasksTemplateRequest struct {
	GoalID            int32  `json:"goal_id" validate:"required,gte=0"`
	Title             string `json:"title" validate:"required,min=3,max=256"`
	ScheduledDatetime string `json:"scheduled_datetime" validate:"required"`
	ScheduledEndTime  string `json:"scheduled_end_time" validate:"omitempty"`
	HasTime           bool   `json:"has_time" validate:"required"`
	RecurrenceRrule   string `json:"recurrence_rrule" validate:"required,min=3"`
}

type RecurringTasksTemplateResponse struct {
	ID                int64  `json:"id"`
	UserID            int32  `json:"user_id"`
	GoalID            int32  `json:"goal_id"`
	Title             string `json:"title"`
	ScheduledDatetime string `json:"scheduled_datetime"`
	HasTime           bool   `json:"has_time"`
	DurationMinutes   int32  `json:"duration_minutes"`
	RecurrenceRrule   string `json:"recurrence_rrule"`
	LastGeneratedDate string `json:"last_generated_date"`
	CreatedAt         string `json:"created_at"`
}

func ToRecurringTasksTemplateResponse(template *domain.RecurringTasksTemplateOutput) RecurringTasksTemplateResponse {
	return RecurringTasksTemplateResponse{
		ID:                template.ID,
		UserID:            template.UserID,
		GoalID:            template.GoalID,
		Title:             template.Title,
		ScheduledDatetime: template.ScheduledDatetime.String(),
		HasTime:           template.HasTime,
		DurationMinutes:   template.DurationMinutes,
		RecurrenceRrule:   template.RecurrenceRrule,
		LastGeneratedDate: template.LastGeneratedDate.String(),
		CreatedAt:         template.CreatedAt.String(),
	}
}

type RecurringTasksTemplateData struct {
	ID                int64  `json:"id"`
	GoalID            int32  `json:"goal_id"`
	Title             string `json:"title"`
	ScheduledDatetime string `json:"scheduled_datetime"`
	HasTime           bool   `json:"has_time"`
	DurationMinutes   int32  `json:"duration_minutes"`
	RecurrenceRrule   string `json:"recurrence_rrule"`
	LastGeneratedDate string `json:"last_generated_date"`
	CreatedAt         string `json:"created_at"`
}

type ListRecurringTasksTemplatesResponse struct {
	UserID                     int32                        `json:"user_id"`
	RecurringTasksTemplateData []RecurringTasksTemplateData `json:"recurring_tasks_template_data"`
}

func ToListRecurringTasksTemplatesResponse(output []domain.RecurringTasksTemplateOutput) ListRecurringTasksTemplatesResponse {
	outTemplatesData := make([]RecurringTasksTemplateData, len(output))

	for index, template := range output {
		outTemplatesData[index] = RecurringTasksTemplateData{
			ID:                template.ID,
			GoalID:            template.GoalID,
			Title:             template.Title,
			ScheduledDatetime: template.ScheduledDatetime.String(),
			HasTime:           template.HasTime,
			DurationMinutes:   template.DurationMinutes,
			RecurrenceRrule:   template.RecurrenceRrule,
			LastGeneratedDate: template.LastGeneratedDate.String(),
			CreatedAt:         template.CreatedAt.String(),
		}
	}

	return ListRecurringTasksTemplatesResponse{
		UserID:                     output[0].UserID,
		RecurringTasksTemplateData: outTemplatesData,
	}
}
