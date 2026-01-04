package domain

import (
	"time"

	repo "github.com/ali-nur31/mile-do/internal/db"
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

type UpdateGoalOutput struct {
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
