package service

import (
	"context"
	"fmt"
	"log/slog"

	repo "github.com/ali-nur31/mile-do/internal/db"
	"github.com/ali-nur31/mile-do/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthTokenManager interface {
	CreateToken(id int64, email string) (string, error)
}

type AuthPasswordManager interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type UserService interface {
	GetUser(ctx context.Context, email string) (repo.User, error)
	CreateUser(ctx context.Context, user domain.UserInput) (domain.UserOutput, error)
	LoginUser(ctx context.Context, user domain.UserInput) (domain.UserOutput, error)
}

type AuthService struct {
	repo            repo.Querier
	tokenManager    AuthTokenManager
	passwordManager AuthPasswordManager
}

func NewUserService(repo repo.Querier, tokenManager AuthTokenManager, passwordManager AuthPasswordManager) UserService {
	return &AuthService{
		repo:            repo,
		tokenManager:    tokenManager,
		passwordManager: passwordManager,
	}
}

func (s *AuthService) GetUser(ctx context.Context, email string) (repo.User, error) {
	return s.repo.GetUser(ctx, email)
}

func (s *AuthService) CreateUser(ctx context.Context, user domain.UserInput) (domain.UserOutput, error) {
	passwordHash, err := s.passwordManager.HashPassword(user.Password)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return domain.UserOutput{}, err
	}

	convertedPasswordHash := pgtype.Text{
		String: passwordHash,
		Valid:  true,
	}

	newUser := repo.CreateUserParams{
		Email:        user.Email,
		PasswordHash: convertedPasswordHash,
	}

	savedUser, err := s.repo.CreateUser(ctx, newUser)
	if err != nil {
		slog.Error("failed to create user", "error", err)
		return domain.UserOutput{}, err
	}

	token, err := s.tokenManager.CreateToken(savedUser.ID, user.Email)

	output := domain.UserOutput{
		Token: token,
	}

	return output, nil
}

func (s *AuthService) LoginUser(ctx context.Context, user domain.UserInput) (domain.UserOutput, error) {
	dbUser, err := s.GetUser(ctx, user.Email)
	if err != nil {
		slog.Error("cannot find user", "error", err)
		return domain.UserOutput{}, err
	}

	passwordIsCorrect := s.passwordManager.CheckPasswordHash(user.Password, dbUser.PasswordHash.String)
	if !passwordIsCorrect {
		slog.Error("password is not correct")
		return domain.UserOutput{}, fmt.Errorf("password is not correct")
	}

	token, _ := s.tokenManager.CreateToken(dbUser.ID, user.Email)

	output := domain.UserOutput{
		Token: token,
	}

	return output, nil
}
