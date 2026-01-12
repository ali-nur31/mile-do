package domain

import (
	"context"
	"time"

	"github.com/ali-nur31/mile-do/internal/repository/db"
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

type TaskService interface {
	ListTasksByGoalID(ctx context.Context, userId int32, goalId int32) ([]TaskOutput, error)
	ListInboxTasks(ctx context.Context, userId int32) ([]TaskOutput, error)
	ListTasksByPeriod(ctx context.Context, period GetTasksByPeriodInput) ([]TaskOutput, error)
	ListTasks(ctx context.Context, userId int32) ([]TaskOutput, error)
	GetTaskByID(ctx context.Context, id int64, userId int32) (*TaskOutput, error)
	CreateTask(ctx context.Context, input CreateTaskInput) (*TaskOutput, error)
	UpdateTask(ctx context.Context, dbTask TaskOutput, updatingTask UpdateTaskInput) (*TaskOutput, error)
	AnalyzeForToday(ctx context.Context, userId int32) (*TodayProgressOutput, error)
	DeleteTaskByID(ctx context.Context, id int64, userId int32) error
	DeleteFutureTasksByRecurringTasksTemplateID(ctx context.Context, templateId int64) error
	CreateTasksByRecurringTasksTemplatesDueForGeneration(ctx context.Context, qtx repo.Querier) error
	CreateTasksByRecurringTasksTemplate(ctx context.Context, qtx repo.Querier, template domain.RecurringTasksTemplateOutput) error
}
