package service

import (
	"context"
	"fmt"
	"time"

	repo "github.com/ali-nur31/mile-do/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateGoalInput struct {
	UserID       int32  `json:"user_id"`
	Title        string `json:"title"`
	Color        string `json:"color"`
	CategoryType string `json:"category_type"`
}

type UpdateGoalInput struct {
	ID           int64  `json:"id"`
	UserID       int32  `json:"user_id"`
	Title        string `json:"title"`
	Color        string `json:"color"`
	CategoryType string `json:"category_type"`
	IsArchived   bool   `json:"is_archived"`
}

type GoalOutput struct {
	ID           int64     `json:"id"`
	UserID       int32     `json:"user_id"`
	Title        string    `json:"title"`
	Color        string    `json:"color"`
	CategoryType string    `json:"category_type"`
	IsArchived   bool      `json:"is_archived"`
	CreatedAt    time.Time `json:"created_at"`
}

type UpdateGoalOutput struct {
	ID           int64  `json:"id"`
	UserID       int32  `json:"user_id"`
	Title        string `json:"title"`
	Color        string `json:"color"`
	CategoryType string `json:"category_type"`
	IsArchived   bool   `json:"is_archived"`
}

type GoalService interface {
	ListGoals(ctx context.Context, filter string) (*[]GoalOutput, error)
	GetGoalByID(ctx context.Context, id int64) (*GoalOutput, error)
	CreateGoal(ctx context.Context, input CreateGoalInput) (*GoalOutput, error)
	UpdateGoal(ctx context.Context, input UpdateGoalInput) (*UpdateGoalOutput, error)
	DeleteGoalByID(ctx context.Context, id int64) error
}

type goalService struct {
	repo repo.Querier
}

func NewGoalService(repo repo.Querier) GoalService {
	return &goalService{
		repo: repo,
	}
}

func (s goalService) ListGoals(ctx context.Context, filter string) (*[]GoalOutput, error) {
	var output []GoalOutput
	var goals []repo.Goal
	var err error

	switch filter {
	case "active":
		goals, err = s.repo.ListGoalsByIsArchived(ctx, false)
	case "archive":
		goals, err = s.repo.ListGoalsByIsArchived(ctx, true)
	case "":
		goals, err = s.repo.ListGoals(ctx)
	default:
		return nil, fmt.Errorf("wrong arguments, expected 'active', 'archive' or ''")
	}

	if err != nil {
		return nil, err
	}

	for _, goal := range goals {
		output = append(output, GoalOutput{
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

func (s goalService) GetGoalByID(ctx context.Context, id int64) (*GoalOutput, error) {
	goal, err := s.repo.GetGoalByID(ctx, id)
	if err != nil {
		return nil, err
	}

	outGoal := GoalOutput{
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

func (s goalService) CreateGoal(ctx context.Context, input CreateGoalInput) (*GoalOutput, error) {
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

	return &GoalOutput{
		ID:           goal.ID,
		UserID:       goal.UserID,
		Title:        goal.Title,
		Color:        goal.Color.String,
		CategoryType: string(goal.CategoryType),
		IsArchived:   goal.IsArchived,
		CreatedAt:    goal.CreatedAt.Time,
	}, nil
}

func (s goalService) UpdateGoal(ctx context.Context, input UpdateGoalInput) (*UpdateGoalOutput, error) {
	goalUpdatingParams := repo.UpdateGoalByIDParams{
		ID:    input.ID,
		Title: input.Title,
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

	return &UpdateGoalOutput{
		ID:           input.ID,
		UserID:       input.UserID,
		Title:        input.Title,
		Color:        input.Color,
		CategoryType: input.CategoryType,
		IsArchived:   input.IsArchived,
	}, nil
}

func (s goalService) DeleteGoalByID(ctx context.Context, id int64) error {
	return s.DeleteGoalByID(ctx, id)
}
