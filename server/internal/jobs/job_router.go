package jobs

import (
	"github.com/ali-nur31/mile-do/config"
	"github.com/ali-nur31/mile-do/internal/domain"
	"github.com/ali-nur31/mile-do/internal/jobs/workers"
	asynq2 "github.com/hibiken/asynq"
)

type JobRouter struct {
	server                        *asynq2.Server
	recurringTasksTemplatesWorker *workers.RecurringTasksTemplatesWorker
}

func NewJobRouter(cfg *config.Redis, recurringTasksTemplatesWorker *workers.RecurringTasksTemplatesWorker) *JobRouter {
	server := asynq2.NewServer(
		asynq2.RedisClientOpt{
			Addr:     cfg.Addr,
			Password: cfg.Password,
			DB:       cfg.DB,
		},
		asynq2.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	return &JobRouter{
		server:                        server,
		recurringTasksTemplatesWorker: recurringTasksTemplatesWorker,
	}
}

func (w *JobRouter) Run() error {
	mux := asynq2.NewServeMux()

	mux.HandleFunc(domain.TypeGenerateRecurringTasks, w.recurringTasksTemplatesWorker.GenerateRecurringTasks)
	mux.HandleFunc(domain.TypeGenerateRecurringTasksByTemplate, w.recurringTasksTemplatesWorker.GenerateRecurringTasksByTemplate)
	mux.HandleFunc(domain.TypeDeleteRecurringTasksByTemplateID, w.recurringTasksTemplatesWorker.DeleteRecurringTasksByTemplateID)

	return w.server.Run(mux)
}
