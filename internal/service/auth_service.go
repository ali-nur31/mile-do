package service

import (
	"context"
	"fmt"

	repo "github.com/ali-nur31/mile-do/internal/db"
	"github.com/ali-nur31/mile-do/internal/domain"
	"github.com/ali-nur31/mile-do/pkg/auth"
)

type AuthTokenManager interface {
	CreateTokens(id int64) (*auth.TokensData, error)
	VerifyToken(tokenString, tokenType string) (*auth.Claims, error)
}

type AuthPasswordManager interface {
	CheckPasswordHash(password, hash string) bool
}

type AuthService interface {
	RegisterUser(ctx context.Context, user domain.UserInput) (*domain.AuthOutput, error)
	LoginUser(ctx context.Context, user domain.UserInput) (*domain.AuthOutput, error)
	LogoutUser(ctx context.Context, userId int32) error
	RefreshTokens(ctx context.Context, refreshToken string) (*domain.AuthOutput, error)
}

type authService struct {
	repo                repo.Querier
	userService         UserService
	tokenManager        AuthTokenManager
	refreshTokenService RefreshTokenService
	passwordManager     AuthPasswordManager
}

func NewAuthService(repo repo.Querier, userService UserService, tokenManager AuthTokenManager, refreshTokenService RefreshTokenService, passwordManager AuthPasswordManager) AuthService {
	return &authService{
		repo:                repo,
		userService:         userService,
		tokenManager:        tokenManager,
		refreshTokenService: refreshTokenService,
		passwordManager:     passwordManager,
	}
}

func (s *authService) RegisterUser(ctx context.Context, user domain.UserInput) (*domain.AuthOutput, error) {
	savedUser, err := s.userService.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	tokensData, err := s.GenerateNewTokens(ctx, savedUser.ID)
	if err != nil {
		return nil, err
	}

	return domain.ToAuthOutput(tokensData), nil
}

func (s *authService) LoginUser(ctx context.Context, user domain.UserInput) (*domain.AuthOutput, error) {
	dbUser, err := s.userService.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	passwordIsCorrect := s.passwordManager.CheckPasswordHash(user.Password, dbUser.PasswordHash)
	if !passwordIsCorrect {
		return nil, fmt.Errorf("password is incorrect")
	}

	tokensData, err := s.GenerateNewTokens(ctx, dbUser.ID)
	if err != nil {
		return nil, err
	}

	return domain.ToAuthOutput(tokensData), nil
}

func (s *authService) LogoutUser(ctx context.Context, userId int32) error {
	err := s.refreshTokenService.DeleteRefreshTokenByUserID(ctx, userId)
	if err != nil {
		return err
	}

	return nil
}

func (s *authService) RefreshTokens(ctx context.Context, refreshToken string) (*domain.AuthOutput, error) {
	claims, err := s.tokenManager.VerifyToken(refreshToken, "refresh")
	if err != nil {
		return nil, err
	}

	dbRefreshToken, err := s.refreshTokenService.GetRefreshTokenByUserID(ctx, int32(claims.ID))
	if err != nil {
		return nil, err
	}

	if dbRefreshToken.TokenHash == "blocked" {
		return nil, fmt.Errorf("user has been banned")
	}

	userId := dbRefreshToken.UserID

	tokensData, err := s.GenerateNewTokens(ctx, int64(userId))
	if err != nil {
		return nil, err
	}

	err = s.refreshTokenService.DeleteRefreshTokenByUserID(ctx, dbRefreshToken.UserID)
	if err != nil {
		return nil, err
	}

	return domain.ToAuthOutput(tokensData), nil
}

func (s *authService) GenerateNewTokens(ctx context.Context, userId int64) (*auth.TokensData, error) {
	tokensData, err := s.tokenManager.CreateTokens(userId)
	if err != nil {
		return nil, err
	}

	err = s.refreshTokenService.CreateRefreshToken(ctx, domain.CreateRefreshTokenInput{
		UserID:    int32(userId),
		TokenHash: tokensData.RefreshToken,
		ExpiresAt: tokensData.RefreshTokenExp,
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't create new refresh token: %w", err)
	}

	return tokensData, nil
}
