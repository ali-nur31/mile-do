package service

import (
	"context"
	"fmt"

	"github.com/ali-nur31/mile-do/internal/domain"
	"github.com/ali-nur31/mile-do/internal/repository/db"
	asynq2 "github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
)

type authService struct {
	repo                repo.Querier
	asynq               *asynq2.Client
	pool                *pgxpool.Pool
	userService         domain.UserService
	tokenManager        domain.AuthTokenManager
	refreshTokenService RefreshTokenService
	passwordManager     domain.AuthPasswordManager
}

func NewAuthService(repo repo.Querier, asynq *asynq2.Client, pool *pgxpool.Pool, userService domain.UserService, tokenManager domain.AuthTokenManager, refreshTokenService RefreshTokenService, passwordManager domain.AuthPasswordManager) domain.AuthService {
	return &authService{
		repo:                repo,
		asynq:               asynq,
		pool:                pool,
		userService:         userService,
		tokenManager:        tokenManager,
		refreshTokenService: refreshTokenService,
		passwordManager:     passwordManager,
	}
}

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

	_, err = s.asynq.Enqueue(domain.NewGenerateDefaultGoalsTask(int32(savedUser.ID)), asynq2.Queue("critical"))
	if err != nil {
		return nil, fmt.Errorf("couldn't enqueue generation of default tasks for new user: %w", err)
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
