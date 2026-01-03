package asynq_jobs

import (
	"github.com/ali-nur31/mile-do/config"
	asynq2 "github.com/hibiken/asynq"
)

type Asynq struct {
	Client *asynq2.Client
}

func InitializeAsynqClient(cfg *config.Redis) (*Asynq, error) {
	asynq := asynq2.NewClient(asynq2.RedisClientOpt{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return &Asynq{Client: asynq}, nil
}
