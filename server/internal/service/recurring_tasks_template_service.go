package service

import (
	"context"
	"fmt"
	"time"

	repo "github.com/ali-nur31/mile-do/internal/db"
	"github.com/ali-nur31/mile-do/internal/domain"
	asynq2 "github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *recurringTasksTemplateService) ListRecurringTasksTemplates(ctx context.Context, userId int32) ([]domain.RecurringTasksTemplateOutput, error) {
	recurringTasksTemplates, err := s.repo.ListRecurringTasksTemplates(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("couldn't get recurring tasks templates: %w", err)
	}

	return domain.ToRecurringTasksTemplateOutputList(recurringTasksTemplates), nil
}

func (s *recurringTasksTemplateService) GetRecurringTasksTemplateByID(ctx context.Context, id int64, userId int32) (*domain.RecurringTasksTemplateOutput, error) {
	template, err := s.repo.GetRecurringTasksTemplateByID(ctx, repo.GetRecurringTasksTemplateByIDParams{
		ID:     id,
		UserID: userId,
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't get recurring tasks template by id: %w", err)
	}

	return domain.ToRecurringTasksTemplateOutput(&template), nil
}

func (s *recurringTasksTemplateService) CreateRecurringTasksTemplate(ctx context.Context, input domain.CreateRecurringTasksTemplateInput) (*domain.RecurringTasksTemplateOutput, error) {
	template, err := s.repo.CreateRecurringTasksTemplate(ctx, repo.CreateRecurringTasksTemplateParams{
		UserID: input.UserID,
		GoalID: input.GoalID,
		Title:  input.Title,
		ScheduledDatetime: pgtype.Timestamp{
			Time:  input.ScheduledDatetime,
			Valid: !input.ScheduledDatetime.IsZero(),
		},
		HasTime:         input.HasTime,
		DurationMinutes: input.DurationMinutes,
		RecurrenceRrule: input.RecurrenceRrule,
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't create new recurring tasks template: %w", err)
	}

	outTemplate := domain.ToRecurringTasksTemplateOutput(&template)

	_, err = s.asynq.Enqueue(domain.NewGenerateRecurringTasksByTemplateTask(outTemplate), asynq2.Queue("critical"))
	if err != nil {
		return nil, fmt.Errorf("couldn't enqueue generation of recurring tasks by template task: %w", err)
	}

	return outTemplate, nil
}

func (s *recurringTasksTemplateService) UpdateRecurringTasksTemplateByID(ctx context.Context, dbTemplate domain.RecurringTasksTemplateOutput, updatingTemplate domain.UpdateRecurringTasksTemplateInput) (*domain.RecurringTasksTemplateOutput, error) {
	templateUpdatingParams := repo.UpdateRecurringTasksTemplateByIDParams{
		ID:     updatingTemplate.ID,
		UserID: updatingTemplate.UserID,
		GoalID: updatingTemplate.GoalID,
		Title:  updatingTemplate.Title,
		ScheduledDatetime: pgtype.Timestamp{
			Time:  updatingTemplate.ScheduledDatetime,
			Valid: !updatingTemplate.ScheduledDatetime.IsZero(),
		},
		HasTime:         updatingTemplate.HasTime,
		DurationMinutes: updatingTemplate.DurationMinutes,
	}

	template, err := s.repo.UpdateRecurringTasksTemplateByID(ctx, templateUpdatingParams)
	if err != nil {
		return nil, fmt.Errorf("couldn't update recurring tasks template: %w", err)
	}

	outTemplate := domain.ToRecurringTasksTemplateOutput(&template)

	_, err = s.asynq.Enqueue(domain.NewDeleteRecurringTasksByTemplateIDTask(dbTemplate.ID), asynq2.Queue("critical"))
	if err != nil {
		return nil, fmt.Errorf("couldn't enqueue deletion of recurring tasks by template id task: %w", err)
	}

	_, err = s.asynq.Enqueue(domain.NewGenerateRecurringTasksByTemplateTask(outTemplate), asynq2.Queue("critical"), asynq2.ProcessIn(1*time.Second))
	if err != nil {
		return nil, fmt.Errorf("couldn't enqueue generation of recurring tasks by template task: %w", err)
	}

	return outTemplate, nil
}

func (s *recurringTasksTemplateService) DeleteRecurringTasksTemplateByID(ctx context.Context, id int64, userId int32) error {
	err := s.repo.DeleteRecurringTasksTemplateByID(ctx, repo.DeleteRecurringTasksTemplateByIDParams{
		ID:     id,
		UserID: userId,
	})
	if err != nil {
		return fmt.Errorf("couldn't delete recurring tasks template by id: %w", err)
	}

	return nil
}

func (s *recurringTasksTemplateService) ListRecurringTasksTemplatesDueForGeneration(ctx context.Context, qtx repo.Querier) ([]domain.RecurringTasksTemplateOutput, error) {
	return s.listRecurringTasksTemplatesDueForGenerationInternal(ctx, qtx)
}

func (s *recurringTasksTemplateService) UpdateLastGeneratedDateInRecurringTasksTemplateByID(ctx context.Context, qtx repo.Querier, updatingTemplate domain.UpdateLastGeneratedDateInRecurringTasksTemplateInput) error {
	templateUpdatingParams := repo.UpdateLastGeneratedDateInRecurringTasksTemplateByIDParams{
		ID: updatingTemplate.ID,
		LastGeneratedDate: pgtype.Date{
			Time:  updatingTemplate.LastGeneratedDate,
			Valid: !updatingTemplate.LastGeneratedDate.IsZero(),
		},
	}

	err := s.updateLastGeneratedDateInRecurringTasksTemplateInternal(ctx, qtx, templateUpdatingParams)
	if err != nil {
		return err
	}

	return nil
}
