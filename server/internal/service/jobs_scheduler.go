package service

import (
	"log/slog"

	"github.com/ali-nur31/mile-do/internal/domain"
	asynq2 "github.com/hibiken/asynq"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron  *cron.Cron
	asynq *asynq2.Client
}

func NewScheduler(cron *cron.Cron, asynq *asynq2.Client) *Scheduler {
	return &Scheduler{
		cron:  cron,
		asynq: asynq,
	}
}

func (s *Scheduler) InitSchedules() {
	s.cron.AddFunc("@daily", func() {
		_, err := s.asynq.Enqueue(domain.NewGenerateRecurringTasksDueForGenerationTask(), asynq2.Queue("default"))
		if err != nil {
			slog.Error("couldn't enqueue generation of recurring tasks due for generation", "error", err)
		}
	})
}
