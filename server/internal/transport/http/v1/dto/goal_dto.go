package dto

import (
	"github.com/ali-nur31/mile-do/internal/domain"
)

type CreateGoalRequest struct {
	Title        string `json:"title" validate:"required,min=3,max=256"`
	Color        string `json:"color" validate:"omitempty,hexcolor"`
	CategoryType string `json:"category_type" validate:"required,oneof=growth maintenance other"`
}

type UpdateGoalRequest struct {
	ID           int64  `json:"id" validate:"required,gte=0"`
	Title        string `json:"title" validate:"required,min=3,max=256"`
	Color        string `json:"color" validate:"omitempty,hexcolor"`
	CategoryType string `json:"category_type" validate:"required,oneof=growth maintenance other"`
	IsArchived   bool   `json:"is_archived" validate:"required,oneof=true false"`
}

type GoalResponse struct {
	ID           int64  `json:"id"`
	UserID       int32  `json:"user_id"`
	Title        string `json:"title"`
	Color        string `json:"color"`
	CategoryType string `json:"category_type"`
	IsArchived   bool   `json:"is_archived"`
	CreatedAt    string `json:"created_at"`
}

func ToGoalResponse(output *domain.GoalOutput) GoalResponse {
	return GoalResponse{
		ID:           output.ID,
		UserID:       output.UserID,
		Title:        output.Title,
		Color:        output.Color,
		CategoryType: output.CategoryType,
		IsArchived:   output.IsArchived,
		CreatedAt:    output.CreatedAt.String(),
	}
}

type GoalData struct {
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	Color        string `json:"color"`
	CategoryType string `json:"category_type"`
	IsArchived   bool   `json:"is_archived"`
	CreatedAt    string `json:"created_at"`
}

type ListGoalsResponse struct {
	UserID int32      `json:"user_id"`
	Data   []GoalData `json:"data"`
}

func ToListGoalsResponse(goals []domain.GoalOutput) ListGoalsResponse {
	outGoalData := make([]GoalData, len(goals))

	for index, goal := range goals {
		outGoalData[index] = GoalData{
			ID:           goal.ID,
			Title:        goal.Title,
			Color:        goal.Color,
			CategoryType: goal.CategoryType,
			IsArchived:   goal.IsArchived,
			CreatedAt:    goal.CreatedAt.String(),
		}
	}

	return ListGoalsResponse{
		UserID: goals[0].UserID,
		Data:   outGoalData,
	}
}
