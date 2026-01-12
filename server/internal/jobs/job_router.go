package jobs

import (
	"github.com/ali-nur31/mile-do/config"
	"github.com/ali-nur31/mile-do/internal/domain"
	"github.com/ali-nur31/mile-do/internal/jobs/workers"
	"github.com/hibiken/asynq"
)

type JobRouter struct {
	server                        *asynq.Server
	goalsWorker                   *workers.GoalsWorker
	recurringTasksTemplatesWorker *workers.RecurringTasksTemplatesWorker
}

func NewJobRouter(
	cfg *config.Redis,
	goalsWorker *workers.GoalsWorker,
	recurringTasksTemplatesWorker *workers.RecurringTasksTemplatesWorker,
) *JobRouter {
	server := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     cfg.Addr,
			Password: cfg.Password,
			DB:       cfg.DB,
		},
		asynq.Config{
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
		goalsWorker:                   goalsWorker,
		recurringTasksTemplatesWorker: recurringTasksTemplatesWorker,
	}
}

func (w *JobRouter) Run() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(domain.TypeGenerateDefaultGoals, w.goalsWorker.GenerateDefaultTasks)

	mux.HandleFunc(domain.TypeGenerateRecurringTasksDueForGeneration, w.recurringTasksTemplatesWorker.GenerateRecurringTasksDueForGeneration)
	mux.HandleFunc(domain.TypeGenerateRecurringTasksByTemplate, w.recurringTasksTemplatesWorker.GenerateRecurringTasksByTemplate)
	mux.HandleFunc(domain.TypeDeleteRecurringTasksByTemplateID, w.recurringTasksTemplatesWorker.DeleteRecurringTasksByTemplateID)

	return w.server.Run(mux)
}
