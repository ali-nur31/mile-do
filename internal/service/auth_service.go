package service

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	repo "github.com/ali-nur31/mile-do/internal/db"
	"github.com/ali-nur31/mile-do/internal/domain"
	"github.com/ali-nur31/mile-do/pkg/auth"
)

type AuthTokenManager interface {
	CreateTokens(id int64) (auth.TokensData, error)
	VerifyRefreshToken(tokenString string) (*auth.RefreshClaims, error)
}

type AuthPasswordManager interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.UserOutput, error)
	GetUserByID(ctx context.Context, id int64) (*domain.UserOutput, error)
	CreateUser(ctx context.Context, user domain.UserInput) (*domain.AuthOutput, error)
	LoginUser(ctx context.Context, user domain.UserInput) (*domain.AuthOutput, error)
	LogoutUser(ctx context.Context, userId int32) error
	RefreshTokens(ctx context.Context, refreshToken string) (*domain.AuthOutput, error)
}

type AuthService struct {
	repo                repo.Querier
	tokenManager        AuthTokenManager
	refreshTokenService RefreshTokenService
	passwordManager     AuthPasswordManager
}

func NewUserService(repo repo.Querier, tokenManager AuthTokenManager, refreshTokenService RefreshTokenService, passwordManager AuthPasswordManager) UserService {
	return &AuthService{
		repo:                repo,
		tokenManager:        tokenManager,
		refreshTokenService: refreshTokenService,
		passwordManager:     passwordManager,
	}
}

func (s *AuthService) GetUserByEmail(ctx context.Context, email string) (*domain.UserOutput, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		slog.Error("couldn't to get user by email: "+email, "error", err)
		return nil, err
	}

	return domain.ToUserOutput(&user), nil
}

func (s *AuthService) GetUserByID(ctx context.Context, id int64) (*domain.UserOutput, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		slog.Error("couldn't get user by id: "+strconv.FormatInt(id, 10), "error", err)
		return nil, err
	}

	return domain.ToUserOutput(&user), nil
}

func (s *AuthService) CreateUser(ctx context.Context, user domain.UserInput) (*domain.AuthOutput, error) {
	passwordHash, err := s.passwordManager.HashPassword(user.Password)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return nil, err
	}

	savedUser, err := s.repo.CreateUser(ctx, repo.CreateUserParams{
		Email:        user.Email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		slog.Error("failed to create user", "error", err)
		return nil, err
	}

	tokensData, err := s.GenerateNewTokens(ctx, savedUser.ID)
	if err != nil {
		slog.Error("failed to generate new tokens", "error", err)
	}

	return domain.ToAuthOutput(tokensData), nil
}

func (s *AuthService) LoginUser(ctx context.Context, user domain.UserInput) (*domain.AuthOutput, error) {
	dbUser, err := s.GetUserByEmail(ctx, user.Email)
	if err != nil {
		slog.Error("cannot find user", "error", err)
		return nil, err
	}

	passwordIsCorrect := s.passwordManager.CheckPasswordHash(user.Password, dbUser.PasswordHash)
	if !passwordIsCorrect {
		slog.Error("password is not correct")
		return nil, fmt.Errorf("password is not correct")
	}

	tokensData, err := s.GenerateNewTokens(ctx, dbUser.ID)
	if err != nil {
		slog.Error("failed to generate new tokens", "error", err)
	}

	return domain.ToAuthOutput(tokensData), nil
}

func (s *AuthService) LogoutUser(ctx context.Context, userId int32) error {
	err := s.refreshTokenService.DeleteRefreshTokenByUserID(ctx, userId)
	if err != nil {
		slog.Error("failed to delete refresh token by user id", "error", err)
		return err
	}

	return nil
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (*domain.AuthOutput, error) {
	claims, err := s.tokenManager.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify refresh token: %v", err)
	}

	dbRefreshToken, err := s.refreshTokenService.GetRefreshTokenByUserID(ctx, int32(claims.ID))
	if err != nil {
		return nil, fmt.Errorf("failed to verify refresh token: %v", err)
	}

	if dbRefreshToken.TokenHash == "blocked" {
		return nil, fmt.Errorf("user has been banned")
	}

	userId := dbRefreshToken.UserID

	tokensData, err := s.GenerateNewTokens(ctx, int64(userId))
	if err != nil {
		slog.Error("failed to generate new tokens", "error", err)
	}

	err = s.refreshTokenService.DeleteRefreshTokenByUserID(ctx, dbRefreshToken.UserID)
	if err != nil {
		return nil, fmt.Errorf("cannot delete old refresh token, err: %v", err)
	}

	return domain.ToAuthOutput(tokensData), nil
}

func (s *AuthService) GenerateNewTokens(ctx context.Context, userId int64) (*auth.TokensData, error) {
	tokensData, err := s.tokenManager.CreateTokens(userId)
	if err != nil {
		slog.Error("failed to generate tokens", "error", err)
		return nil, err
	}

	err = s.refreshTokenService.CreateRefreshToken(ctx, domain.CreateRefreshTokenInput{
		UserID:    int32(userId),
		TokenHash: tokensData.RefreshToken,
		ExpiresAt: tokensData.RefreshTokenExp,
	})
	if err != nil {
		slog.Error("failed create new refresh token", "error", err)
		return nil, err
	}

	return &tokensData, nil
}
