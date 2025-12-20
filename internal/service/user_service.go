package service

import (
	"context"

	repo "github.com/ali-nur31/mile-do/internal/db"
)

type UserService interface {
	GetUser(ctx context.Context, email string) (repo.User, error)
}

type svc struct {
	repo repo.Querier
}

func NewUserService(repo repo.Querier) UserService {
	return &svc{
		repo: repo,
	}
}

func (s *svc) GetUser(ctx context.Context, email string) (repo.User, error) {
	return s.repo.GetUser(ctx, email)
}
