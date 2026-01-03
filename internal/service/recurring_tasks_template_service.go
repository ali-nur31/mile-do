package service

import (
	"context"
	"time"

	repo "github.com/ali-nur31/mile-do/internal/db"
	"github.com/ali-nur31/mile-do/internal/domain"
	asynq2 "github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
)

type RecurringTasksTemplateService interface {
	ListRecurringTasksTemplates(ctx context.Context, userId int32) ([]domain.RecurringTasksTemplateOutput, error)
	GetRecurringTasksTemplateByID(ctx context.Context, id int64, userId int32) (*domain.RecurringTasksTemplateOutput, error)
	CreateRecurringTasksTemplate(ctx context.Context, input domain.CreateRecurringTasksTemplateInput) (*domain.RecurringTasksTemplateOutput, error)
	UpdateRecurringTasksTemplateByID(ctx context.Context, dbTemplate domain.RecurringTasksTemplateOutput, updatingTemplate domain.UpdateRecurringTasksTemplateInput) (*domain.RecurringTasksTemplateOutput, error)
	DeleteRecurringTasksTemplateByID(ctx context.Context, id int64, userId int32) error
	ListRecurringTasksTemplatesDueForGeneration(ctx context.Context) ([]domain.RecurringTasksTemplateOutput, error)
	UpdateLastGeneratedDateInRecurringTasksTemplateByID(ctx context.Context, updatingTemplate domain.UpdateLastGeneratedDateInRecurringTasksTemplateInput) (*domain.LastGeneratedDate, error)
}

type recurringTasksTemplateService struct {
	repo  repo.Querier
	asynq *asynq2.Client
}

func NewRecurringTasksTemplateService(repo repo.Querier, asynq *asynq2.Client) RecurringTasksTemplateService {
	return &recurringTasksTemplateService{
		repo:  repo,
		asynq: asynq,
	}
}

func (s *recurringTasksTemplateService) ListRecurringTasksTemplates(ctx context.Context, userId int32) ([]domain.RecurringTasksTemplateOutput, error) {
	recurringTasksTemplates, err := s.repo.ListRecurringTasksTemplates(ctx, userId)
	if err != nil {
		return nil, err
	}

	output := make([]domain.RecurringTasksTemplateOutput, 0, len(recurringTasksTemplates))

	for _, template := range recurringTasksTemplates {
		output = append(output, domain.RecurringTasksTemplateOutput{
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
		})
	}

	return output, nil
}

func (s *recurringTasksTemplateService) GetRecurringTasksTemplateByID(ctx context.Context, id int64, userId int32) (*domain.RecurringTasksTemplateOutput, error) {
	template, err := s.repo.GetRecurringTasksTemplateByID(ctx, repo.GetRecurringTasksTemplateByIDParams{
		ID:     id,
		UserID: userId,
	})
	if err != nil {
		return nil, err
	}

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

	return &outTemplate, nil
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
		return nil, err
	}

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

	_, err = s.asynq.Enqueue(domain.NewGenerateRecurringTasksByTemplateTask(outTemplate), asynq2.Queue("critical"))
	if err != nil {
		return nil, err
	}

	return &outTemplate, nil
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

	err := s.repo.UpdateRecurringTasksTemplateByID(ctx, templateUpdatingParams)
	if err != nil {
		return nil, err
	}

	outTemplate := domain.RecurringTasksTemplateOutput{
		ID:                templateUpdatingParams.ID,
		UserID:            templateUpdatingParams.UserID,
		GoalID:            templateUpdatingParams.GoalID,
		Title:             templateUpdatingParams.Title,
		ScheduledDatetime: templateUpdatingParams.ScheduledDatetime.Time,
		HasTime:           templateUpdatingParams.HasTime,
		DurationMinutes:   templateUpdatingParams.DurationMinutes,
		RecurrenceRrule:   templateUpdatingParams.RecurrenceRrule,
		CreatedAt:         dbTemplate.CreatedAt,
	}

	_, err = s.asynq.Enqueue(domain.NewDeleteRecurringTasksByTemplateIDTask(dbTemplate.ID), asynq2.Queue("critical"))
	if err != nil {
		return nil, err
	}

	_, err = s.asynq.Enqueue(domain.NewGenerateRecurringTasksByTemplateTask(outTemplate), asynq2.Queue("critical"), asynq2.ProcessIn(1*time.Second))
	if err != nil {
		return nil, err
	}

	return &outTemplate, nil
}

func (s *recurringTasksTemplateService) DeleteRecurringTasksTemplateByID(ctx context.Context, id int64, userId int32) error {
	return s.repo.DeleteRecurringTasksTemplateByID(ctx, repo.DeleteRecurringTasksTemplateByIDParams{
		ID:     id,
		UserID: userId,
	})
}

func (s *recurringTasksTemplateService) ListRecurringTasksTemplatesDueForGeneration(ctx context.Context) ([]domain.RecurringTasksTemplateOutput, error) {
	recurringTasksTemplates, err := s.repo.ListRecurringTasksTemplatesDueForGeneration(ctx)
	if err != nil {
		return nil, err
	}

	output := make([]domain.RecurringTasksTemplateOutput, 0, len(recurringTasksTemplates))

	for _, template := range recurringTasksTemplates {
		output = append(output, domain.RecurringTasksTemplateOutput{
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
		})
	}

	return output, nil
}

func (s *recurringTasksTemplateService) UpdateLastGeneratedDateInRecurringTasksTemplateByID(ctx context.Context, updatingTemplate domain.UpdateLastGeneratedDateInRecurringTasksTemplateInput) (*domain.LastGeneratedDate, error) {
	templateUpdatingParams := repo.UpdateLastGeneratedDateInRecurringTasksTemplateByIDParams{
		ID: updatingTemplate.ID,
		LastGeneratedDate: pgtype.Date{
			Time:  updatingTemplate.LastGeneratedDate,
			Valid: !updatingTemplate.LastGeneratedDate.IsZero(),
		},
	}

	err := s.repo.UpdateLastGeneratedDateInRecurringTasksTemplateByID(ctx, templateUpdatingParams)
	if err != nil {
		return nil, err
	}

	lastGeneratedDate := domain.LastGeneratedDate(updatingTemplate.LastGeneratedDate)

	return &lastGeneratedDate, nil
}
