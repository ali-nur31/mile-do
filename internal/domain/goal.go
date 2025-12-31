package domain

import "time"

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
