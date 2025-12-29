package domain

import "time"

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

type UpdateGoalOutput struct {
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
