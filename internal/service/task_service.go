package service

import (
	"context"
	"fmt"
	"log/slog"

	repo "github.com/ali-nur31/mile-do/internal/db"
	"github.com/ali-nur31/mile-do/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *taskService) ListTasksByGoalID(ctx context.Context, userId int32, goalId int32) ([]domain.TaskOutput, error) {
	tasks, err := s.repo.ListTasksByGoalID(ctx, repo.ListTasksByGoalIDParams{
		UserID: userId,
		GoalID: goalId,
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't get tasks by goal id: %w", err)
	}

	return domain.ToTaskOutputList(tasks), nil
}

func (s *taskService) ListInboxTasks(ctx context.Context, userId int32) ([]domain.TaskOutput, error) {
	tasks, err := s.repo.ListInboxTasks(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("couldn't get inbox tasks: %w", err)
	}

	return domain.ToTaskOutputList(tasks), nil
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
		return nil, fmt.Errorf("couldn't get tasks by period: %w", err)
	}

	return domain.ToTaskOutputList(tasks), nil
}

func (s *taskService) ListTasks(ctx context.Context, userId int32) ([]domain.TaskOutput, error) {
	tasks, err := s.repo.ListTasks(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("couldn't get tasks: %w", err)
	}

	return domain.ToTaskOutputList(tasks), nil
}

func (s *taskService) GetTaskByID(ctx context.Context, id int64, userId int32) (*domain.TaskOutput, error) {
	task, err := s.repo.GetTaskByID(ctx, repo.GetTaskByIDParams{
		ID:     id,
		UserID: userId,
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't get task by id: %w", err)
	}

	return domain.ToTaskOutput(&task), nil
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
			Microseconds: convertTimeToMicroseconds(input.ScheduledTime),
			Valid:        input.HasTime,
		},
		DurationMinutes: pgtype.Int4{
			Int32: int32(input.DurationMinutes.Minutes()),
			Valid: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't create task: %w", err)
	}

	return domain.ToTaskOutput(&task), nil
}

func (s *taskService) UpdateTask(ctx context.Context, dbTask domain.TaskOutput, updatingTask domain.UpdateTaskInput) (*domain.TaskOutput, error) {
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
			Microseconds: convertTimeToMicroseconds(updatingTask.ScheduledTime),
			Valid:        updatingTask.HasTime,
		},
		DurationMinutes: pgtype.Int4{
			Int32: int32(updatingTask.DurationMinutes.Minutes()),
			Valid: true,
		},
	}

	task, err := s.repo.UpdateTaskByID(ctx, taskUpdatingParams)
	if err != nil {
		return nil, fmt.Errorf("couldn't update task: %w", err)
	}

	return domain.ToTaskOutput(&task), nil
}

func (s *taskService) AnalyzeForToday(ctx context.Context, userId int32) (*domain.TodayProgressOutput, error) {
	stats, err := s.repo.CountCompletedTasksForToday(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("couldn't get tasks statistics for today: %w", err)
	}

	return &domain.TodayProgressOutput{
		TotalTasks:     stats.TotalToday,
		CompletedToday: stats.CompletedToday,
	}, nil
}

func (s *taskService) DeleteTaskByID(ctx context.Context, id int64, userId int32) error {
	err := s.repo.DeleteTaskByID(ctx, repo.DeleteTaskByIDParams{
		ID:     id,
		UserID: userId,
	})
	if err != nil {
		return fmt.Errorf("couldn't delete task by id: %w", err)
	}

	return nil
}

func (s *taskService) DeleteFutureTasksByRecurringTasksTemplateID(ctx context.Context, templateId int64) error {
	err := s.repo.DeleteFutureTasksByRecurringTasksTemplateID(ctx, pgtype.Int4{
		Int32: int32(templateId),
		Valid: true,
	})
	if err != nil {
		return fmt.Errorf("couldn't delete future tasks by recurring tasks template id: %w", err)
	}

	return nil
}

func (s *taskService) CreateTasksByRecurringTasksTemplates(ctx context.Context) error {
	templates, err := s.repo.ListRecurringTasksTemplatesDueForGeneration(ctx)
	if err != nil {
		return err
	}

	for _, template := range templates {
		outTemplate := domain.RecurringTasksTemplateOutput{
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

		err = s.CreateTasksByRecurringTasksTemplate(ctx, outTemplate)
		if err != nil {
			slog.Error("failed to process template", "template_id", template.ID, "error", err)
			continue
		}
	}

	return nil
}

func (s *taskService) CreateTasksByRecurringTasksTemplate(ctx context.Context, template domain.RecurringTasksTemplateOutput) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(context.Background())
	}()

	qtx := repo.New(tx)

	err = s.CreateTasksByTemplateInternal(ctx, template, qtx)

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("couldn't commit transaction for create tasks by recurring tasks template: %w", err)
	}

	return nil
}
