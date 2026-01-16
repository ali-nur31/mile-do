package service

import (
	"context"
	"fmt"

	"github.com/ali-nur31/mile-do/internal/domain"
	repo "github.com/ali-nur31/mile-do/internal/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *authService) generateNewTokensInternal(ctx context.Context, qtx repo.Querier, userId int64) (*domain.TokensData, error) {
	tokensData, err := s.tokenManager.CreateTokens(userId)
	if err != nil {
		return nil, err
	}

	_, err = qtx.CreateRefreshToken(ctx, repo.CreateRefreshTokenParams{
		UserID:    int32(userId),
		TokenHash: tokensData.RefreshToken,
		ExpiresAt: pgtype.Timestamp{
			Time:  tokensData.RefreshTokenExp,
			Valid: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't create new refresh token: %w", err)
	}

	return tokensData, nil
}
