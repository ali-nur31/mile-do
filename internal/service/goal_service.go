package service

import (
	"context"
	"fmt"

	repo "github.com/ali-nur31/mile-do/internal/db"
	"github.com/ali-nur31/mile-do/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

type GoalService interface {
	ListGoals(ctx context.Context, filter string, userId int32) ([]domain.GoalOutput, error)
	GetGoalByID(ctx context.Context, id int64, userId int32) (*domain.GoalOutput, error)
	CreateGoal(ctx context.Context, input domain.CreateGoalInput) (*domain.GoalOutput, error)
	UpdateGoal(ctx context.Context, input domain.UpdateGoalInput) (*domain.GoalOutput, error)
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

func (s *goalService) ListGoals(ctx context.Context, filter string, userId int32) ([]domain.GoalOutput, error) {
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
		return nil, fmt.Errorf("wrong arguments, expected 'active', 'archive' or nothing")
	}

	if err != nil {
		return nil, fmt.Errorf("couldn't get goals: %w", err)
	}

	return domain.ToGoalOutputList(goals), nil
}

func (s *goalService) GetGoalByID(ctx context.Context, id int64, userId int32) (*domain.GoalOutput, error) {
	goal, err := s.repo.GetGoalByID(ctx, repo.GetGoalByIDParams{
		ID:     id,
		UserID: userId,
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't get goal by id: %w", err)
	}

	return domain.ToGoalOutput(&goal), nil
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
		return nil, fmt.Errorf("couldn't create new goal: %w", err)
	}

	return domain.ToGoalOutput(&goal), nil
}

func (s *goalService) UpdateGoal(ctx context.Context, input domain.UpdateGoalInput) (*domain.GoalOutput, error) {
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

	goal, err := s.repo.UpdateGoalByID(ctx, goalUpdatingParams)
	if err != nil {
		return nil, fmt.Errorf("couldn't update goal: %w", err)
	}

	return domain.ToGoalOutput(&goal), nil
}

func (s *goalService) DeleteGoalByID(ctx context.Context, id int64, userId int32) error {
	err := s.repo.DeleteGoalByID(ctx, repo.DeleteGoalByIDParams{
		ID:     id,
		UserID: userId,
	})
	if err != nil {
		return fmt.Errorf("couldn't delete goal by id: %w", err)
	}

	return nil
}
