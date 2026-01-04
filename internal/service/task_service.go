package service

import (
	"context"
	"log/slog"
	"time"

	repo "github.com/ali-nur31/mile-do/internal/db"
	"github.com/ali-nur31/mile-do/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/teambition/rrule-go"
)

type TaskService interface {
	ListTasksByGoalID(ctx context.Context, userId int32, goalId int32) ([]domain.TaskOutput, error)
	ListInboxTasks(ctx context.Context, userId int32) ([]domain.TaskOutput, error)
	ListTasksByPeriod(ctx context.Context, period domain.GetTasksByPeriodInput) ([]domain.TaskOutput, error)
	ListTasks(ctx context.Context, userId int32) ([]domain.TaskOutput, error)
	GetTaskByID(ctx context.Context, id int64, userId int32) (*domain.TaskOutput, error)
	CreateTask(ctx context.Context, input domain.CreateTaskInput) (*domain.TaskOutput, error)
	UpdateTask(ctx context.Context, dbTask domain.TaskOutput, updatingTask domain.UpdateTask) (*domain.TaskOutput, error)
	AnalyzeForToday(ctx context.Context, userId int32) (repo.CountCompletedTasksForTodayRow, error)
	DeleteTaskByID(ctx context.Context, id int64, userId int32) error
	DeleteFutureTasksByRecurringTasksTemplateID(ctx context.Context, templateId int64) error
	CreateTasksByRecurringTasksTemplates(ctx context.Context) error
	CreateTasksByRecurringTasksTemplate(ctx context.Context, template domain.RecurringTasksTemplateOutput) error
}

type taskService struct {
	repo                          repo.Querier
	recurringTasksTemplateService RecurringTasksTemplateService
}

func NewTaskService(repo repo.Querier, recurringTasksTemplateService RecurringTasksTemplateService) TaskService {
	return &taskService{
		repo:                          repo,
		recurringTasksTemplateService: recurringTasksTemplateService,
	}
}

func (s *taskService) ListTasksByGoalID(ctx context.Context, userId int32, goalId int32) ([]domain.TaskOutput, error) {
	tasks, err := s.repo.ListTasksByGoalID(ctx, repo.ListTasksByGoalIDParams{
		UserID: userId,
		GoalID: goalId,
	})
	if err != nil {
		return nil, err
	}

	return mapTasksToOutputList(tasks), nil
}

func (s *taskService) ListInboxTasks(ctx context.Context, userId int32) ([]domain.TaskOutput, error) {
	tasks, err := s.repo.ListInboxTasks(ctx, userId)
	if err != nil {
		return nil, err
	}

	return mapTasksToOutputList(tasks), nil
}

func (s *taskService) ListTasksByPeriod(ctx context.Context, period domain.GetTasksByPeriodInput) ([]domain.TaskOutput, error) {
	tasks, err := s.repo.ListTasksByDateRange(ctx, repo.ListTasksByDateRangeParams{
		UserID: period.UserID,
		ScheduledDate: pgtype.Date{
			Time:  period.AfterDate,
			Valid: true,
		},
		ScheduledDate_2: pgtype.Date{
			Time:  period.BeforeDate,
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}

	return mapTasksToOutputList(tasks), nil
}

func (s *taskService) ListTasks(ctx context.Context, userId int32) ([]domain.TaskOutput, error) {
	tasks, err := s.repo.ListTasks(ctx, userId)
	if err != nil {
		return nil, err
	}

	return mapTasksToOutputList(tasks), nil
}

func (s *taskService) GetTaskByID(ctx context.Context, id int64, userId int32) (*domain.TaskOutput, error) {
	task, err := s.repo.GetTaskByID(ctx, repo.GetTaskByIDParams{
		ID:     id,
		UserID: userId,
	})
	if err != nil {
		return nil, err
	}

	outTask := domain.TaskOutput{
		ID:              task.ID,
		UserID:          task.UserID,
		GoalID:          task.GoalID,
		Title:           task.Title,
		IsDone:          task.IsDone,
		ScheduledDate:   task.ScheduledDate.Time,
		ScheduledTime:   microsecondsToTime(task.ScheduledTime.Microseconds),
		HasTime:         task.HasTime,
		RescheduleCount: task.RescheduleCount,
		CreatedAt:       task.CreatedAt.Time,
	}

	return &outTask, nil
}

func (s *taskService) CreateTask(ctx context.Context, input domain.CreateTaskInput) (*domain.TaskOutput, error) {
	task, err := s.repo.CreateTask(ctx, repo.CreateTaskParams{
		UserID: input.UserID,
		GoalID: input.GoalID,
		Title:  input.Title,
		ScheduledDate: pgtype.Date{
			Time:  input.ScheduledDate,
			Valid: !input.ScheduledDate.IsZero(),
		},
		HasTime: input.HasTime,
		ScheduledTime: pgtype.Time{
			Microseconds: timeToMicroseconds(input.ScheduledTime),
			Valid:        input.HasTime,
		},
		DurationMinutes: pgtype.Int4{
			Int32: int32(input.DurationMinutes.Minutes()),
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}

	return &domain.TaskOutput{
		ID:              task.ID,
		UserID:          task.UserID,
		GoalID:          task.GoalID,
		Title:           task.Title,
		IsDone:          task.IsDone,
		ScheduledDate:   task.ScheduledDate.Time,
		ScheduledTime:   microsecondsToTime(task.ScheduledTime.Microseconds),
		HasTime:         task.HasTime,
		DurationMinutes: task.DurationMinutes.Int32,
		RescheduleCount: task.RescheduleCount,
		CreatedAt:       task.CreatedAt.Time,
	}, nil
}

func (s *taskService) UpdateTask(ctx context.Context, dbTask domain.TaskOutput, updatingTask domain.UpdateTask) (*domain.TaskOutput, error) {
	if !dbTask.ScheduledDate.IsZero() && !dbTask.ScheduledDate.Equal(updatingTask.ScheduledDate) {
		updatingTask.RescheduleCount += 1
	}

	taskUpdatingParams := repo.UpdateTaskByIDParams{
		ID:     updatingTask.ID,
		UserID: updatingTask.UserID,
		GoalID: updatingTask.GoalID,
		Title:  updatingTask.Title,
		IsDone: updatingTask.IsDone,
		ScheduledDate: pgtype.Date{
			Time:  updatingTask.ScheduledDate,
			Valid: !updatingTask.ScheduledDate.IsZero(),
		},
		HasTime: updatingTask.HasTime,
		ScheduledTime: pgtype.Time{
			Microseconds: timeToMicroseconds(updatingTask.ScheduledTime),
			Valid:        updatingTask.HasTime,
		},
		DurationMinutes: pgtype.Int4{
			Int32: int32(updatingTask.DurationMinutes.Minutes()),
			Valid: true,
		},
	}

	err := s.repo.UpdateTaskByID(ctx, taskUpdatingParams)
	if err != nil {
		return nil, err
	}

	return &domain.TaskOutput{
		ID:              taskUpdatingParams.ID,
		UserID:          taskUpdatingParams.UserID,
		GoalID:          taskUpdatingParams.GoalID,
		Title:           taskUpdatingParams.Title,
		IsDone:          taskUpdatingParams.IsDone,
		ScheduledDate:   taskUpdatingParams.ScheduledDate.Time,
		ScheduledTime:   microsecondsToTime(taskUpdatingParams.ScheduledTime.Microseconds),
		HasTime:         taskUpdatingParams.HasTime,
		DurationMinutes: taskUpdatingParams.DurationMinutes.Int32,
		RescheduleCount: taskUpdatingParams.RescheduleCount,
		CreatedAt:       dbTask.CreatedAt,
	}, nil
}

func (s *taskService) AnalyzeForToday(ctx context.Context, userId int32) (repo.CountCompletedTasksForTodayRow, error) {
	return s.repo.CountCompletedTasksForToday(ctx, userId)
}

func (s *taskService) DeleteTaskByID(ctx context.Context, id int64, userId int32) error {
	return s.repo.DeleteTaskByID(ctx, repo.DeleteTaskByIDParams{
		ID:     id,
		UserID: userId,
	})
}

func (s *taskService) DeleteFutureTasksByRecurringTasksTemplateID(ctx context.Context, templateId int64) error {
	return s.repo.DeleteFutureTasksByRecurringTasksTemplateID(ctx, pgtype.Int4{
		Int32: int32(templateId),
		Valid: true,
	})
}

func (s *taskService) CreateTasksByRecurringTasksTemplates(ctx context.Context) error {
	var templates []domain.RecurringTasksTemplateOutput
	var err error

	templates, err = s.recurringTasksTemplateService.ListRecurringTasksTemplatesDueForGeneration(ctx)
	if err != nil {
		return err
	}

	horizonDate := time.Now().UTC().AddDate(0, 1, 0)
	for _, template := range templates {
		var rule *rrule.Set

		rule, err = rrule.StrToRRuleSet(template.RecurrenceRrule)
		if err != nil {
			slog.Error("failed to parse rrule from template", "error", err)
			return err
		}

		rule.DTStart(template.ScheduledDatetime)

		var lastDate time.Time
		if template.LastGeneratedDate.IsZero() {
			lastDate = template.ScheduledDatetime.Add(-1 * time.Second)
		} else {
			lastDate = template.LastGeneratedDate
		}

		dates := rule.Between(lastDate, horizonDate, true)
		if len(dates) == 0 {
			continue
		}

		for _, date := range dates {
			var scheduledDateOnly, scheduledTimeOnly time.Time

			scheduledDateOnly, _ = time.Parse(time.DateOnly, date.Format(time.DateOnly))
			scheduledTimeOnly, _ = time.Parse(time.TimeOnly, date.Format(time.TimeOnly))

			_, err = s.CreateTask(ctx, domain.CreateTaskInput{
				UserID:          template.UserID,
				GoalID:          template.GoalID,
				Title:           template.Title,
				ScheduledDate:   scheduledDateOnly,
				ScheduledTime:   scheduledTimeOnly,
				HasTime:         template.HasTime,
				DurationMinutes: time.Duration(template.DurationMinutes) * time.Minute,
			})
			if err != nil {
				slog.Error("cannot create task by recurring template", "error", err)
				return err
			}
		}

		newLastGeneratedDate := dates[len(dates)-1]

		_, err = s.recurringTasksTemplateService.UpdateLastGeneratedDateInRecurringTasksTemplateByID(ctx, domain.UpdateLastGeneratedDateInRecurringTasksTemplateInput{
			ID:                template.ID,
			LastGeneratedDate: newLastGeneratedDate,
		})
		if err != nil {
			slog.Error("cannot update last_generated_date in recurring_template", "error", err)
			return err
		}
	}

	return nil
}

func (s *taskService) CreateTasksByRecurringTasksTemplate(ctx context.Context, template domain.RecurringTasksTemplateOutput) error {
	var err error

	horizonDate := time.Now().UTC().AddDate(0, 3, 0)
	var rule *rrule.Set

	rule, err = rrule.StrToRRuleSet(template.RecurrenceRrule)
	if err != nil {
		slog.Error("failed to parse rrule from template", "error", err)
		return err
	}

	rule.DTStart(template.ScheduledDatetime)

	var lastDate time.Time
	if template.LastGeneratedDate.IsZero() {
		lastDate = template.ScheduledDatetime.Add(-1 * time.Second)
	} else {
		lastDate = template.LastGeneratedDate
	}

	dates := rule.Between(lastDate, horizonDate, true)
	if len(dates) == 0 {
		return nil
	}

	for _, date := range dates {
		var scheduledDateOnly, scheduledTimeOnly time.Time

		scheduledDateOnly, _ = time.Parse(time.DateOnly, date.Format(time.DateOnly))
		scheduledTimeOnly, _ = time.Parse(time.TimeOnly, date.Format(time.TimeOnly))

		_, err = s.CreateTask(ctx, domain.CreateTaskInput{
			UserID:          template.UserID,
			GoalID:          template.GoalID,
			Title:           template.Title,
			ScheduledDate:   scheduledDateOnly,
			ScheduledTime:   scheduledTimeOnly,
			HasTime:         template.HasTime,
			DurationMinutes: time.Duration(template.DurationMinutes) * time.Minute,
		})
		if err != nil {
			slog.Error("cannot create task by recurring template", "error", err)
			return err
		}
	}

	newLastGeneratedDate := dates[len(dates)-1]

	_, err = s.recurringTasksTemplateService.UpdateLastGeneratedDateInRecurringTasksTemplateByID(ctx, domain.UpdateLastGeneratedDateInRecurringTasksTemplateInput{
		ID:                template.ID,
		LastGeneratedDate: newLastGeneratedDate,
	})
	if err != nil {
		slog.Error("cannot update last_generated_date in recurring_template", "error", err)
		return err
	}

	return nil
}

func mapTasksToOutputList(tasks []repo.Task) []domain.TaskOutput {
	output := make([]domain.TaskOutput, 0, len(tasks))

	for _, task := range tasks {
		output = append(output, domain.TaskOutput{
			ID:              task.ID,
			UserID:          task.UserID,
			GoalID:          task.GoalID,
			Title:           task.Title,
			IsDone:          task.IsDone,
			ScheduledDate:   task.ScheduledDate.Time,
			ScheduledTime:   microsecondsToTime(task.ScheduledTime.Microseconds),
			HasTime:         task.HasTime,
			DurationMinutes: task.DurationMinutes.Int32,
			RescheduleCount: task.RescheduleCount,
			CreatedAt:       task.CreatedAt.Time,
		})
	}

	return output
}

func microsecondsToTime(msec int64) time.Time {
	return time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(msec) * time.Microsecond)
}

func timeToMicroseconds(t time.Time) int64 {
	return int64(t.Hour())*3600000000 +
		int64(t.Minute())*60000000 +
		int64(t.Second())*1000000 +
		int64(t.Nanosecond())/1000
}
