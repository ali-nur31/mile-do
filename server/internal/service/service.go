package service

import (
	"context"

	repo "github.com/ali-nur31/mile-do/internal/db"
	"github.com/ali-nur31/mile-do/internal/domain"
	"github.com/ali-nur31/mile-do/pkg/auth"
	asynq2 "github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthTokenManager interface {
	CreateTokens(id int64) (*auth.TokensData, error)
	VerifyToken(tokenString, tokenType string) (*auth.Claims, error)
}

type AuthPasswordManager interface {
	HashPassword(password string) (string, error)
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
	asynq               *asynq2.Client
	pool                *pgxpool.Pool
	userService         UserService
	tokenManager        AuthTokenManager
	refreshTokenService RefreshTokenService
	passwordManager     AuthPasswordManager
}

func NewAuthService(repo repo.Querier, asynq *asynq2.Client, pool *pgxpool.Pool, userService UserService, tokenManager AuthTokenManager, refreshTokenService RefreshTokenService, passwordManager AuthPasswordManager) AuthService {
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

type RefreshTokenService interface {
	GetRefreshTokenByUserID(ctx context.Context, qtx repo.Querier, userId int32) (*domain.RefreshTokenOutput, error)
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

type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.UserOutput, error)
	GetUserByID(ctx context.Context, id int64) (*domain.UserOutput, error)
	CreateUser(ctx context.Context, qtx repo.Querier, user domain.UserInput) (*domain.UserOutput, error)
}

type userService struct {
	repo            repo.Querier
	passwordManager AuthPasswordManager
}

func NewUserService(repo repo.Querier, passwordManager AuthPasswordManager) UserService {
	return &userService{
		repo:            repo,
		passwordManager: passwordManager,
	}
}

type GoalService interface {
	ListGoals(ctx context.Context, filter string, userId int32) ([]domain.GoalOutput, error)
	GetGoalByID(ctx context.Context, id int64, userId int32) (*domain.GoalOutput, error)
	CreateGoal(ctx context.Context, qtx repo.Querier, input domain.CreateGoalInput) (*domain.GoalOutput, error)
	UpdateGoal(ctx context.Context, input domain.UpdateGoalInput) (*domain.GoalOutput, error)
	DeleteGoalByID(ctx context.Context, id int64, userId int32) error
}

type goalService struct {
	repo repo.Querier
}

func NewGoalService(repo repo.Querier) GoalService {
	return &goalService{
		repo: repo,
	}
}

type RecurringTasksTemplateService interface {
	ListRecurringTasksTemplates(ctx context.Context, userId int32) ([]domain.RecurringTasksTemplateOutput, error)
	GetRecurringTasksTemplateByID(ctx context.Context, id int64, userId int32) (*domain.RecurringTasksTemplateOutput, error)
	CreateRecurringTasksTemplate(ctx context.Context, input domain.CreateRecurringTasksTemplateInput) (*domain.RecurringTasksTemplateOutput, error)
	UpdateRecurringTasksTemplateByID(ctx context.Context, dbTemplate domain.RecurringTasksTemplateOutput, updatingTemplate domain.UpdateRecurringTasksTemplateInput) (*domain.RecurringTasksTemplateOutput, error)
	DeleteRecurringTasksTemplateByID(ctx context.Context, id int64, userId int32) error
	ListRecurringTasksTemplatesDueForGeneration(ctx context.Context) ([]domain.RecurringTasksTemplateOutput, error)
	UpdateLastGeneratedDateInRecurringTasksTemplateByID(ctx context.Context, updatingTemplate domain.UpdateLastGeneratedDateInRecurringTasksTemplateInput) (string, error)
}

type recurringTasksTemplateService struct {
	repo  repo.Querier
	asynq *asynq2.Client
}

func NewRecurringTasksTemplateService(repo repo.Querier, asynq *asynq2.Client) RecurringTasksTemplateService {
	return &recurringTasksTemplateService{
		repo:  repo,
		asynq: asynq,
	}
}

type TaskService interface {
	ListTasksByGoalID(ctx context.Context, userId int32, goalId int32) ([]domain.TaskOutput, error)
	ListInboxTasks(ctx context.Context, userId int32) ([]domain.TaskOutput, error)
	ListTasksByPeriod(ctx context.Context, period domain.GetTasksByPeriodInput) ([]domain.TaskOutput, error)
	ListTasks(ctx context.Context, userId int32) ([]domain.TaskOutput, error)
	GetTaskByID(ctx context.Context, id int64, userId int32) (*domain.TaskOutput, error)
	CreateTask(ctx context.Context, input domain.CreateTaskInput) (*domain.TaskOutput, error)
	UpdateTask(ctx context.Context, dbTask domain.TaskOutput, updatingTask domain.UpdateTaskInput) (*domain.TaskOutput, error)
	AnalyzeForToday(ctx context.Context, userId int32) (*domain.TodayProgressOutput, error)
	DeleteTaskByID(ctx context.Context, id int64, userId int32) error
	DeleteFutureTasksByRecurringTasksTemplateID(ctx context.Context, templateId int64) error
	CreateTasksByRecurringTasksTemplates(ctx context.Context) error
	CreateTasksByRecurringTasksTemplate(ctx context.Context, template domain.RecurringTasksTemplateOutput) error
}

type taskService struct {
	repo                          repo.Querier
	pool                          *pgxpool.Pool
	recurringTasksTemplateService RecurringTasksTemplateService
}

func NewTaskService(repo repo.Querier, pool *pgxpool.Pool, recurringTasksTemplateService RecurringTasksTemplateService) TaskService {
	return &taskService{
		repo:                          repo,
		pool:                          pool,
		recurringTasksTemplateService: recurringTasksTemplateService,
	}
}
