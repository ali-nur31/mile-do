package asynq_jobs

import (
	"log/slog"

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

	if err := asynq.Ping(); err != nil {
		slog.Error("failed to connect to Redis through Asynq client", "error", err)
		return nil, err
	}

	return &Asynq{Client: asynq}, nil
}
