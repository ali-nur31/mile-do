package service

import (
	"context"
	"log/slog"

	repo "github.com/ali-nur31/mile-do/internal/db"
	"github.com/ali-nur31/mile-do/pkg/auth"
	"github.com/ali-nur31/mile-do/pkg/bcrypt"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserInput struct {
	Email    string
	Password string
}

type UserOutput struct {
	Token string
}

type UserService interface {
	GetUser(ctx context.Context, email string) (repo.User, error)
	CreateUser(ctx context.Context, user UserInput) (UserOutput, error)
	LoginUser(ctx context.Context, user UserInput) (UserOutput, error)
}

type AuthService struct {
	repo         repo.Querier
	tokenManager auth.JwtManager
}

func NewUserService(repo repo.Querier, tokenManager auth.JwtManager) UserService {
	return &AuthService{
		repo:         repo,
		tokenManager: tokenManager,
	}
}

func (s *AuthService) GetUser(ctx context.Context, email string) (repo.User, error) {
	return s.repo.GetUser(ctx, email)
}

func (s *AuthService) CreateUser(ctx context.Context, user UserInput) (UserOutput, error) {
	passwordHash, err := bcrypt.HashPassword(user.Password)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return UserOutput{}, err
	}

	convertedPasswordHash := pgtype.Text{
		String: passwordHash,
	}

	newUser := repo.CreateUserParams{
		Email:        user.Email,
		PasswordHash: convertedPasswordHash,
	}

	_, err = s.repo.CreateUser(ctx, newUser)
	if err != nil {
		slog.Error("failed to create user", "error", err)
		return UserOutput{}, err
	}

	token, err := s.tokenManager.CreateToken(user.Email)

	output := UserOutput{
		Token: token,
	}

	return output, nil
}

func (s *AuthService) LoginUser(ctx context.Context, user UserInput) (UserOutput, error) {
	//dbUser, err := s.GetUser(ctx, user.Email)
	//if err != nil {
	//	slog.Error("cannot find user", "error", err)
	//	return UserOutput{}, err
	//}

	//passwordIsCorrect := bcrypt.CheckPasswordHash(user.Password, dbUser.PasswordHash.String)
	//if !passwordIsCorrect {
	//	slog.Error("password is not correct")
	//	return UserOutput{}, fmt.Errorf("password is not correct")
	//}

	token, _ := s.tokenManager.CreateToken(user.Email)

	output := UserOutput{
		Token: token,
	}

	return output, nil
}
