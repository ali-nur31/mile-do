package domain

import (
	"time"
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

type LastGeneratedDate time.Time

type RecurringTasksTemplateOutput struct {
	ID                int64
	UserID            int32
	GoalID            int32
	Title             string
	ScheduledDatetime time.Time
	HasTime           bool
	DurationMinutes   int32
	RecurrenceRrule   string
	// ToDo change type to LastGeneratedDate
	LastGeneratedDate time.Time
	CreatedAt         time.Time
}
