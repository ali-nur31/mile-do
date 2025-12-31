package domain

import "time"

type GetTasksByPeriodInput struct {
	UserID     int32
	AfterDate  time.Time
	BeforeDate time.Time
}

type CreateTaskInput struct {
	UserID          int32
	GoalID          int32
	Title           string
	ScheduledDate   time.Time
	ScheduledTime   time.Time
	HasTime         bool
	DurationMinutes time.Duration
}

type UpdateTask struct {
	ID              int64
	UserID          int32
	GoalID          int32
	Title           string
	IsDone          bool
	ScheduledDate   time.Time
	ScheduledTime   time.Time
	HasTime         bool
	DurationMinutes time.Duration
	RescheduleCount int32
}

type TaskOutput struct {
	ID              int64
	UserID          int32
	GoalID          int32
	Title           string
	IsDone          bool
	ScheduledDate   time.Time
	ScheduledTime   time.Time
	DurationMinutes int32
	RescheduleCount int32
	CreatedAt       time.Time
}
