package service

import (
	"context"
	"fmt"

	"github.com/ali-nur31/mile-do/internal/domain"
	"github.com/ali-nur31/mile-do/internal/repository/db"
)

type userService struct {
	repo            repo.Querier
	passwordManager domain.AuthPasswordManager
}

func NewUserService(repo repo.Querier, passwordManager domain.AuthPasswordManager) domain.UserService {
	return &userService{
		repo:            repo,
		passwordManager: passwordManager,
	}
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*domain.UserOutput, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("couldn't get user by email: %w", err)
	}

	return domain.ToUserOutput(&user), nil
}

func (s *userService) GetUserByID(ctx context.Context, id int64) (*domain.UserOutput, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("couldn't get user by id: %w", err)
	}

	return domain.ToUserOutput(&user), nil
}

func (s *userService) CreateUser(ctx context.Context, qtx repo.Querier, user domain.UserInput) (*domain.UserOutput, error) {
	return s.createUserInternal(ctx, qtx, user)
}
