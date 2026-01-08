package service

import (
	"context"
	"fmt"
	"time"

	repo "github.com/ali-nur31/mile-do/internal/db"
	"github.com/ali-nur31/mile-do/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/teambition/rrule-go"
)

func (s *taskService) CreateTasksByTemplateInternal(ctx context.Context, template domain.RecurringTasksTemplateOutput, qtx *repo.Queries) error {
	var err error

	horizonDate := time.Now().UTC().AddDate(0, 3, 0)
	var rule *rrule.Set

	rule, err = rrule.StrToRRuleSet(template.RecurrenceRrule)
	if err != nil {
		return fmt.Errorf("couldn't parse rrule from template: %w", err)
	}

	rule.DTStart(template.ScheduledDatetime)

	var lastDate time.Time
	if template.LastGeneratedDate.IsZero() {
		lastDate = template.ScheduledDatetime.Add(-1 * time.Second)
	} else {
		lastDate = template.LastGeneratedDate
	}

	dates := rule.Between(lastDate, horizonDate, true)
	if len(dates) == 0 {
		return nil
	}

	for _, date := range dates {
		scheduledDateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

		_, err = qtx.CreateTask(ctx, repo.CreateTaskParams{
			UserID: template.UserID,
			GoalID: template.GoalID,
			RecurringTemplateID: pgtype.Int4{
				Int32: int32(template.ID),
				Valid: true,
			},
			Title: template.Title,
			ScheduledDate: pgtype.Date{
				Time:  scheduledDateOnly,
				Valid: true,
			},
			ScheduledTime: pgtype.Time{
				Microseconds: convertTimeToMicroseconds(date),
				Valid:        template.HasTime,
			},
			HasTime: template.HasTime,
			DurationMinutes: pgtype.Int4{
				Int32: template.DurationMinutes,
				Valid: true,
			},
		})
		if err != nil {
			return fmt.Errorf("couldn't create task by recurring tasks template: %w", err)
		}
	}

	newLastGeneratedDate := dates[len(dates)-1]

	err = qtx.UpdateLastGeneratedDateInRecurringTasksTemplateByID(ctx, repo.UpdateLastGeneratedDateInRecurringTasksTemplateByIDParams{
		ID: template.ID,
		LastGeneratedDate: pgtype.Date{
			Time:  newLastGeneratedDate,
			Valid: true,
		},
	})
	if err != nil {
		return fmt.Errorf("couldn't update last_generated_date in recurring_template: %w", err)
	}

	return nil
}

func convertTimeToMicroseconds(t time.Time) int64 {
	return int64(t.Hour())*3600000000 +
		int64(t.Minute())*60000000 +
		int64(t.Second())*1000000 +
		int64(t.Nanosecond())/1000
}
