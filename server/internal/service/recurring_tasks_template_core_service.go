package service

import (
	"context"
	"fmt"

	repo "github.com/ali-nur31/mile-do/internal/db"
	"github.com/ali-nur31/mile-do/internal/domain"
)

func (s *recurringTasksTemplateService) listRecurringTasksTemplatesDueForGenerationInternal(ctx context.Context, qtx repo.Querier) ([]domain.RecurringTasksTemplateOutput, error) {
	recurringTasksTemplates, err := qtx.ListRecurringTasksTemplatesDueForGeneration(ctx)
	if err != nil {
		return nil, fmt.Errorf("couldn't get recurring tasks templates due for generation: %w", err)
	}

	output := domain.ToRecurringTasksTemplateOutputList(recurringTasksTemplates)

	return output, nil
}

func (s *recurringTasksTemplateService) updateLastGeneratedDateInRecurringTasksTemplateInternal(ctx context.Context, qtx repo.Querier, updatingTemplate repo.UpdateLastGeneratedDateInRecurringTasksTemplateByIDParams) error {
	err := qtx.UpdateLastGeneratedDateInRecurringTasksTemplateByID(ctx, updatingTemplate)
	if err != nil {
		return fmt.Errorf("couldn't update last_generated_date in recurring_template: %w", err)
	}

	return nil
}
