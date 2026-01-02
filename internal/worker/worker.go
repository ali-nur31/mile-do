package worker

import (
	"github.com/ali-nur31/mile-do/config"
	"github.com/ali-nur31/mile-do/internal/worker/jobs"
	asynq2 "github.com/hibiken/asynq"
)

type Worker struct {
	server                   *asynq2.Server
	taskGenerateRecurringJob *jobs.TaskGenerateRecurringJob
}

func NewWorker(cfg *config.Redis, taskGenerateRecurringJob *jobs.TaskGenerateRecurringJob) *Worker {
	server := asynq2.NewServer(
		asynq2.RedisClientOpt{
			Addr:     cfg.Addr,
			Password: cfg.Password,
			DB:       cfg.DB,
		},
		asynq2.Config{
			Concurrency: 10,
		},
	)

	return &Worker{
		server:                   server,
		taskGenerateRecurringJob: taskGenerateRecurringJob,
	}
}

func (w *Worker) Run() error {
	mux := asynq2.NewServeMux()

	return w.server.Run(mux)
}
