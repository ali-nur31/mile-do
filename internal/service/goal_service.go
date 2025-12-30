package service

import (
	"context"
	"fmt"

	repo "github.com/ali-nur31/mile-do/internal/db"
	"github.com/ali-nur31/mile-do/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

type GoalService interface {
	ListGoals(ctx context.Context, filter string, userId int32) (*[]domain.GoalOutput, error)
	GetGoalByID(ctx context.Context, id int64, userId int32) (*domain.GoalOutput, error)
	CreateGoal(ctx context.Context, input domain.CreateGoalInput) (*domain.GoalOutput, error)
	UpdateGoal(ctx context.Context, input domain.UpdateGoalInput) (*domain.UpdateGoalOutput, error)
	DeleteGoalByID(ctx context.Context, id int64, userId int32) error
}

type goalService struct {
	repo repo.Querier
}

func NewGoalService(repo repo.Querier) GoalService {
	return &goalService{
		repo: repo,
	}
}

func (s *goalService) ListGoals(ctx context.Context, filter string, userId int32) (*[]domain.GoalOutput, error) {
	var output []domain.GoalOutput
	var goals []repo.Goal
	var err error

	switch filter {
	case "active":
		goals, err = s.repo.ListGoalsByIsArchived(ctx, repo.ListGoalsByIsArchivedParams{
			IsArchived: false,
			UserID:     userId,
		})
	case "archive":
		goals, err = s.repo.ListGoalsByIsArchived(ctx, repo.ListGoalsByIsArchivedParams{
			IsArchived: true,
			UserID:     userId,
		})
	case "":
		goals, err = s.repo.ListGoals(ctx, userId)
	default:
		return nil, fmt.Errorf("wrong arguments, expected 'active', 'archive' or ''")
	}

	if err != nil {
		return nil, err
	}

	for _, goal := range goals {
		output = append(output, domain.GoalOutput{
			ID:           goal.ID,
			UserID:       goal.UserID,
			Title:        goal.Title,
			Color:        goal.Color.String,
			CategoryType: string(goal.CategoryType),
			IsArchived:   goal.IsArchived,
			CreatedAt:    goal.CreatedAt.Time,
		})
	}

	return &output, nil
}

func (s *goalService) GetGoalByID(ctx context.Context, id int64, userId int32) (*domain.GoalOutput, error) {
	goal, err := s.repo.GetGoalByID(ctx, repo.GetGoalByIDParams{
		ID:     id,
		UserID: userId,
	})
	if err != nil {
		return nil, err
	}

	outGoal := domain.GoalOutput{
		ID:           goal.ID,
		UserID:       goal.UserID,
		Title:        goal.Title,
		Color:        goal.Color.String,
		CategoryType: string(goal.CategoryType),
		IsArchived:   goal.IsArchived,
		CreatedAt:    goal.CreatedAt.Time,
	}

	return &outGoal, nil
}

func (s *goalService) CreateGoal(ctx context.Context, input domain.CreateGoalInput) (*domain.GoalOutput, error) {
	goal, err := s.repo.CreateGoal(ctx, repo.CreateGoalParams{
		UserID: input.UserID,
		Title:  input.Title,
		Color: pgtype.Text{
			String: input.Color,
			Valid:  true,
		},
		CategoryType: repo.GoalsCategoryType(input.CategoryType),
	})
	if err != nil {
		return nil, err
	}

	return &domain.GoalOutput{
		ID:           goal.ID,
		UserID:       goal.UserID,
		Title:        goal.Title,
		Color:        goal.Color.String,
		CategoryType: string(goal.CategoryType),
		IsArchived:   goal.IsArchived,
		CreatedAt:    goal.CreatedAt.Time,
	}, nil
}

func (s *goalService) UpdateGoal(ctx context.Context, input domain.UpdateGoalInput) (*domain.UpdateGoalOutput, error) {
	goalUpdatingParams := repo.UpdateGoalByIDParams{
		ID:     input.ID,
		UserID: input.UserID,
		Title:  input.Title,
		Color: pgtype.Text{
			String: input.Color,
			Valid:  true,
		},
		CategoryType: repo.GoalsCategoryType(input.CategoryType),
		IsArchived:   input.IsArchived,
	}

	err := s.repo.UpdateGoalByID(ctx, goalUpdatingParams)
	if err != nil {
		return nil, err
	}

	return &domain.UpdateGoalOutput{
		ID:           input.ID,
		UserID:       input.UserID,
		Title:        input.Title,
		Color:        input.Color,
		CategoryType: input.CategoryType,
		IsArchived:   input.IsArchived,
	}, nil
}

func (s *goalService) DeleteGoalByID(ctx context.Context, id int64, userId int32) error {
	return s.repo.DeleteGoalByID(ctx, repo.DeleteGoalByIDParams{
		ID:     id,
		UserID: userId,
	})
}
