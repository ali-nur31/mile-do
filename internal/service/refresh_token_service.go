package service

import (
	"context"
	"fmt"

	repo "github.com/ali-nur31/mile-do/internal/db"
	"github.com/ali-nur31/mile-do/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

type RefreshTokenService interface {
	GetRefreshTokenByUserID(ctx context.Context, userId int32) (repo.RefreshToken, error)
	CreateRefreshToken(ctx context.Context, input domain.CreateRefreshTokenInput) error
	DeleteRefreshTokenByUserID(ctx context.Context, userId int32) error
}

type refreshTokenService struct {
	repo repo.Querier
}

func NewRefreshTokenService(repo repo.Querier) RefreshTokenService {
	return &refreshTokenService{
		repo: repo,
	}
}

func (s *refreshTokenService) GetRefreshTokenByUserID(ctx context.Context, userId int32) (repo.RefreshToken, error) {
	return s.repo.GetRefreshTokenByUserID(ctx, userId)
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
		return fmt.Errorf("failed to save refresh_token: %v", err)
	}

	return nil
}

func (s *refreshTokenService) DeleteRefreshTokenByUserID(ctx context.Context, userId int32) error {
	err := s.repo.DeleteRefreshTokenByUserID(ctx, userId)
	if err != nil {
		return fmt.Errorf("failed to delete refresh_token from db: %v", err)
	}

	return nil
}
