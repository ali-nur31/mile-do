package service

import (
	"context"
	"fmt"
	"log/slog"

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
	GetUser(ctx context.Context, email string) (repo.User, error)
	CreateUser(ctx context.Context, user domain.UserInput) (domain.UserOutput, error)
	LoginUser(ctx context.Context, user domain.UserInput) (domain.UserOutput, error)
	LogoutUser(ctx context.Context, userId int32) error
	RefreshTokens(ctx context.Context, refreshToken string) (domain.UserOutput, error)
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

func (s *AuthService) GetUser(ctx context.Context, email string) (repo.User, error) {
	return s.repo.GetUser(ctx, email)
}

func (s *AuthService) CreateUser(ctx context.Context, user domain.UserInput) (domain.UserOutput, error) {
	passwordHash, err := s.passwordManager.HashPassword(user.Password)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return domain.UserOutput{}, err
	}

	newUser := repo.CreateUserParams{
		Email:        user.Email,
		PasswordHash: passwordHash,
	}

	savedUser, err := s.repo.CreateUser(ctx, newUser)
	if err != nil {
		slog.Error("failed to create user", "error", err)
		return domain.UserOutput{}, err
	}

	tokensData, err := s.tokenManager.CreateTokens(savedUser.ID)
	if err != nil {
		slog.Error("failed to generate tokens", "error", err)
		return domain.UserOutput{}, err
	}

	err = s.refreshTokenService.CreateRefreshToken(ctx, domain.CreateRefreshTokenInput{
		UserID:    int32(savedUser.ID),
		TokenHash: tokensData.RefreshToken,
		ExpiresAt: tokensData.RefreshTokenExp,
	})

	output := domain.UserOutput{
		AccessToken:  tokensData.AccessToken,
		RefreshToken: tokensData.RefreshToken,
	}

	return output, nil
}

func (s *AuthService) LoginUser(ctx context.Context, user domain.UserInput) (domain.UserOutput, error) {
	dbUser, err := s.GetUser(ctx, user.Email)
	if err != nil {
		slog.Error("cannot find user", "error", err)
		return domain.UserOutput{}, err
	}

	passwordIsCorrect := s.passwordManager.CheckPasswordHash(user.Password, dbUser.PasswordHash)
	if !passwordIsCorrect {
		slog.Error("password is not correct")
		return domain.UserOutput{}, fmt.Errorf("password is not correct")
	}

	tokensData, err := s.tokenManager.CreateTokens(dbUser.ID)
	if err != nil {
		slog.Error("failed to generate tokens", "error", err)
		return domain.UserOutput{}, err
	}

	err = s.refreshTokenService.CreateRefreshToken(ctx, domain.CreateRefreshTokenInput{
		UserID:    int32(dbUser.ID),
		TokenHash: tokensData.RefreshToken,
		ExpiresAt: tokensData.RefreshTokenExp,
	})

	output := domain.UserOutput{
		AccessToken:  tokensData.AccessToken,
		RefreshToken: tokensData.RefreshToken,
	}

	return output, nil
}

func (s *AuthService) LogoutUser(ctx context.Context, userId int32) error {
	return s.refreshTokenService.DeleteRefreshTokenByUserID(ctx, userId)
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (domain.UserOutput, error) {
	claims, err := s.tokenManager.VerifyRefreshToken(refreshToken)
	if err != nil {
		return domain.UserOutput{}, fmt.Errorf("failed to verify refresh token: %v", err)
	}

	dbRefreshToken, err := s.refreshTokenService.GetRefreshTokenByUserID(ctx, int32(claims.ID))
	if err != nil {
		return domain.UserOutput{}, fmt.Errorf("failed to verify refresh token: %v", err)
	}

	if dbRefreshToken.TokenHash == "blocked" {
		return domain.UserOutput{}, fmt.Errorf("user has been banned")
	}

	userId := dbRefreshToken.UserID

	tokensData, err := s.tokenManager.CreateTokens(int64(userId))
	if err != nil {
		slog.Error("failed to generate tokens", "error", err)
		return domain.UserOutput{}, err
	}

	err = s.refreshTokenService.DeleteRefreshTokenByUserID(ctx, dbRefreshToken.UserID)
	if err != nil {
		return domain.UserOutput{}, fmt.Errorf("cannot delete old refresh token, err: %v", err)
	}

	err = s.refreshTokenService.CreateRefreshToken(ctx, domain.CreateRefreshTokenInput{
		UserID:    userId,
		TokenHash: tokensData.RefreshToken,
		ExpiresAt: tokensData.RefreshTokenExp,
	})
	if err != nil {
		return domain.UserOutput{}, fmt.Errorf("failed to save refresh token, err: %v", err)
	}

	return domain.UserOutput{
		AccessToken:  tokensData.AccessToken,
		RefreshToken: tokensData.RefreshToken,
	}, nil
}
