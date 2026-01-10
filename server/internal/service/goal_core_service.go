package service

import (
	"context"
	"fmt"

	repo "github.com/ali-nur31/mile-do/internal/db"
	"github.com/ali-nur31/mile-do/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *goalService) createGoalInternal(ctx context.Context, qtx repo.Querier, input domain.CreateGoalInput) (*domain.GoalOutput, error) {
	goal, err := qtx.CreateGoal(ctx, repo.CreateGoalParams{
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
