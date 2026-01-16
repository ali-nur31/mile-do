package service

import (
	"context"
	"fmt"

	repo "github.com/ali-nur31/mile-do/internal/db"
	"github.com/ali-nur31/mile-do/internal/domain"
)

func (s *authService) RegisterUser(ctx context.Context, user domain.UserInput) (*domain.AuthOutput, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(context.Background())
	}()

	qtx := repo.New(tx)

	savedUser, err := s.userService.CreateUser(ctx, qtx, user)
	if err != nil {
		return nil, err
	}

	defaultGoals := []domain.CreateGoalInput{
		{
			UserID:       int32(savedUser.ID),
			Title:        "Routine",
			Color:        "#73260A",
			CategoryType: "maintenance",
		},
		{
			UserID:       int32(savedUser.ID),
			Title:        "Other",
			Color:        "#0096ff",
			CategoryType: "other",
		},
	}

	for _, input := range defaultGoals {
		_, err = s.goalService.CreateGoal(ctx, qtx, input)
		if err != nil {
			return nil, err
		}
	}

	tokensData, err := s.generateNewTokensInternal(ctx, qtx, savedUser.ID)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("couldn't commit transaction for registering user: %w", err)
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

	tokensData, err := s.generateNewTokensInternal(ctx, s.repo, dbUser.ID)
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

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(context.Background())
	}()

	qtx := repo.New(tx)

	dbRefreshToken, err := s.refreshTokenService.GetRefreshTokenByUserID(ctx, qtx, int32(claims.ID))
	if err != nil {
		return nil, err
	}

	if dbRefreshToken.TokenHash == "blocked" {
		return nil, fmt.Errorf("user has been banned")
	}

	userId := dbRefreshToken.UserID

	tokensData, err := s.generateNewTokensInternal(ctx, qtx, int64(userId))
	if err != nil {
		return nil, err
	}

	err = qtx.DeleteRefreshTokenByUserID(ctx, dbRefreshToken.UserID)
	if err != nil {
		return nil, fmt.Errorf("couldn't delete refresh token by user id for refresh tokens: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("couldn't commit transaction for refresh tokens: %w", err)
	}

	return domain.ToAuthOutput(tokensData), nil
}
