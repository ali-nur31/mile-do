package service

import (
	"context"
	"fmt"

	repo "github.com/ali-nur31/mile-do/internal/db"
	"github.com/ali-nur31/mile-do/internal/domain"
)

type UserPasswordManager interface {
	HashPassword(password string) (string, error)
}

type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.UserOutput, error)
	GetUserByID(ctx context.Context, id int64) (*domain.UserOutput, error)
	CreateUser(ctx context.Context, user domain.UserInput) (*domain.UserOutput, error)
}

type userService struct {
	repo            repo.Querier
	passwordManager UserPasswordManager
}

func NewUserService(repo repo.Querier, passwordManager UserPasswordManager) UserService {
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

func (s *userService) CreateUser(ctx context.Context, user domain.UserInput) (*domain.UserOutput, error) {
	passwordHash, err := s.passwordManager.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed when hashing password: %w", err)
	}

	savedUser, err := s.repo.CreateUser(ctx, repo.CreateUserParams{
		Email:        user.Email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't create new user: %w", err)
	}

	return domain.ToUserOutput(&savedUser), nil
}
