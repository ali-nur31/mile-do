package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ali-nur31/mile-do/internal/domain"
	repo "github.com/ali-nur31/mile-do/internal/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type taskService struct {
	repo                          repo.Querier
	pool                          *pgxpool.Pool
	recurringTasksTemplateService domain.RecurringTasksTemplateService
}

func NewTaskService(repo repo.Querier, pool *pgxpool.Pool, recurringTasksTemplateService domain.RecurringTasksTemplateService) domain.TaskService {
	return &taskService{
		repo:                          repo,
		pool:                          pool,
		recurringTasksTemplateService: recurringTasksTemplateService,
	}
}

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
			Int32: input.DurationMinutes,
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
			Int32: updatingTask.DurationMinutes,
			Valid: true,
		},
	}

	task, err := s.repo.UpdateTaskByID(ctx, taskUpdatingParams)
	if err != nil {
		return nil, fmt.Errorf("couldn't update task: %w", err)
	}

	return domain.ToTaskOutput(&task), nil
}

func (s *taskService) CompleteTask(ctx context.Context, userId int32, taskId int64) (*domain.TaskOutput, error) {

	taskUpdatingParams := repo.UpdateIsDoneInTaskByIDParams{
		ID:     taskId,
		UserID: userId,
		IsDone: true,
	}

	task, err := s.repo.UpdateIsDoneInTaskByID(ctx, taskUpdatingParams)
	if err != nil {
		return nil, fmt.Errorf("couldn't complete task: %w", err)
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

func (s *taskService) CreateTasksByRecurringTasksTemplatesDueForGeneration(ctx context.Context, qtx repo.Querier) error {
	templates, err := s.recurringTasksTemplateService.ListRecurringTasksTemplatesDueForGeneration(ctx, qtx)
	if err != nil {
		return err
	}

	if templates == nil {
		return nil
	}

	for _, template := range templates {
		outTemplate := domain.RecurringTasksTemplateOutput{
			ID:                template.ID,
			UserID:            template.UserID,
			GoalID:            template.GoalID,
			Title:             template.Title,
			ScheduledDatetime: template.ScheduledDatetime,
			HasTime:           template.HasTime,
			DurationMinutes:   template.DurationMinutes,
			RecurrenceRrule:   template.RecurrenceRrule,
			LastGeneratedDate: template.LastGeneratedDate,
			CreatedAt:         template.CreatedAt,
		}

		err = s.CreateTasksByRecurringTasksTemplate(ctx, qtx, outTemplate)
		if err != nil {
			slog.Error("failed to process template", "template_id", template.ID, "error", err)
			continue
		}
	}

	return nil
}

func (s *taskService) CreateTasksByRecurringTasksTemplate(ctx context.Context, qtx repo.Querier, template domain.RecurringTasksTemplateOutput) error {
	err := s.CreateTasksByTemplateInternal(ctx, template, qtx)
	if err != nil {
		return err
	}

	return nil
}
