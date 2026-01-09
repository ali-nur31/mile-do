package domain

import (
	"time"

	repo "github.com/ali-nur31/mile-do/internal/db"
)

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
	DurationMinutes int32
}

type UpdateTaskInput struct {
	ID              int64
	UserID          int32
	GoalID          int32
	Title           string
	IsDone          bool
	ScheduledDate   time.Time
	ScheduledTime   time.Time
	HasTime         bool
	DurationMinutes int32
	RescheduleCount int32
}

type TodayProgressOutput struct {
	TotalTasks     int32
	CompletedToday int32
}

type TaskOutput struct {
	ID              int64
	UserID          int32
	GoalID          int32
	Title           string
	IsDone          bool
	ScheduledDate   time.Time
	ScheduledTime   time.Time
	HasTime         bool
	DurationMinutes int32
	RescheduleCount int32
	CreatedAt       time.Time
}

func ToTaskOutput(t *repo.Task) *TaskOutput {
	return &TaskOutput{
		ID:              t.ID,
		UserID:          t.UserID,
		GoalID:          t.GoalID,
		Title:           t.Title,
		IsDone:          t.IsDone,
		ScheduledDate:   t.ScheduledDate.Time,
		ScheduledTime:   microsecondsToTime(t.ScheduledTime.Microseconds),
		HasTime:         t.HasTime,
		DurationMinutes: t.DurationMinutes.Int32,
		CreatedAt:       t.CreatedAt.Time,
	}
}

func ToTaskOutputList(tasks []repo.Task) []TaskOutput {
	output := make([]TaskOutput, len(tasks))
	for i, t := range tasks {
		output[i] = *ToTaskOutput(&t)
	}
	return output
}

func microsecondsToTime(msec int64) time.Time {
	return time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(msec) * time.Microsecond)
}
