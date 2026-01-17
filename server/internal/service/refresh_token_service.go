package service

import (
	"context"
	"fmt"

	"github.com/ali-nur31/mile-do/internal/domain"
	repo "github.com/ali-nur31/mile-do/internal/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type refreshTokenService struct {
	repo repo.Querier
}

func NewRefreshTokenService(repo repo.Querier) domain.RefreshTokenService {
	return &refreshTokenService{
		repo: repo,
	}
}

func (s *refreshTokenService) GetRefreshTokenByUserID(ctx context.Context, qtx repo.Querier, userId int32) (*domain.RefreshTokenOutput, error) {
	return s.getRefreshTokenByUserIDInternal(ctx, qtx, userId)
}

func (s *refreshTokenService) CreateRefreshToken(ctx context.Context, input domain.CreateRefreshTokenInput) error {
	_, err := s.repo.CreateRefreshToken(ctx, repo.CreateRefreshTokenParams{
		UserID:    input.UserID,
		TokenHash: input.TokenHash,
		ExpiresAt: pgtype.Timestamp{
			Time:  input.ExpiresAt,
			Valid: true,
		},
	})
	if err != nil {
		return fmt.Errorf("couldn't save new refresh token: %w", err)
	}

	return nil
}

func (s *refreshTokenService) DeleteRefreshTokenByUserID(ctx context.Context, userId int32) error {
	err := s.repo.DeleteRefreshTokenByUserID(ctx, userId)
	if err != nil {
		return fmt.Errorf("couldn't delete refresh token from db: %w", err)
	}

	return nil
}
