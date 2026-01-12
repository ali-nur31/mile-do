package domain

import (
	"context"
	"time"

	"github.com/ali-nur31/mile-do/internal/repository/db"
)

type CreateGoalInput struct {
	UserID       int32
	Title        string
	Color        string
	CategoryType string
}

type UpdateGoalInput struct {
	ID           int64
	UserID       int32
	Title        string
	Color        string
	CategoryType string
	IsArchived   bool
}

type GoalOutput struct {
	ID           int64
	UserID       int32
	Title        string
	Color        string
	CategoryType string
	IsArchived   bool
	CreatedAt    time.Time
}

func ToGoalOutput(goal *repo.Goal) *GoalOutput {
	return &GoalOutput{
		ID:           goal.ID,
		UserID:       goal.UserID,
		Title:        goal.Title,
		Color:        goal.Color.String,
		CategoryType: string(goal.CategoryType),
		IsArchived:   goal.IsArchived,
		CreatedAt:    goal.CreatedAt.Time,
	}
}

func ToGoalOutputList(goals []repo.Goal) []GoalOutput {
	output := make([]GoalOutput, len(goals))
	for i, g := range goals {
		output[i] = *ToGoalOutput(&g)
	}
	return output
}

type GoalService interface {
	ListGoals(ctx context.Context, filter string, userId int32) ([]GoalOutput, error)
	GetGoalByID(ctx context.Context, id int64, userId int32) (*GoalOutput, error)
	CreateGoal(ctx context.Context, qtx repo.Querier, input CreateGoalInput) (*GoalOutput, error)
	UpdateGoal(ctx context.Context, input UpdateGoalInput) (*GoalOutput, error)
	DeleteGoalByID(ctx context.Context, id int64, userId int32) error
}
