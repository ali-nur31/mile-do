package service

import (
	"context"
	"fmt"

	"github.com/ali-nur31/mile-do/internal/domain"
	repo "github.com/ali-nur31/mile-do/internal/repository/db"
)

func (s *userService) createUserInternal(ctx context.Context, qtx repo.Querier, user domain.UserInput) (*domain.UserOutput, error) {
	passwordHash, err := s.passwordManager.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed when hashing password: %w", err)
	}

	savedUser, err := qtx.CreateUser(ctx, repo.CreateUserParams{
		Email:        user.Email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't create new user: %w", err)
	}

	return domain.ToUserOutput(&savedUser), nil
}
