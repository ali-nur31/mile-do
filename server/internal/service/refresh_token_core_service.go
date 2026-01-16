package service

import (
	"context"
	"fmt"

	"github.com/ali-nur31/mile-do/internal/domain"
	"github.com/ali-nur31/mile-do/internal/repository/db"
)

func (s *refreshTokenService) getRefreshTokenByUserIDInternal(ctx context.Context, qtx repo.Querier, userId int32) (*domain.RefreshTokenOutput, error) {
	refreshToken, err := qtx.GetRefreshTokenByUserID(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("couldn't get refresh token by user id: %w", err)
	}

	return domain.ToRefreshTokenOutput(&refreshToken), nil
}
