package asynq_jobs

import (
	"github.com/ali-nur31/mile-do/config"
	"github.com/hibiken/asynq"
)

type Asynq struct {
	Client *asynq.Client
}

func InitializeAsynqClient(cfg *config.Redis) (*Asynq, error) {
	asynq := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return &Asynq{Client: asynq}, nil
}
