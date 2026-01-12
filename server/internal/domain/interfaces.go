package domain

import (
	"context"
	"time"

	repo "github.com/ali-nur31/mile-do/internal/repository/db"
	"github.com/ali-nur31/mile-do/pkg/auth"
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
	RegisterUser(ctx context.Context, user UserInput) (*AuthOutput, error)
	LoginUser(ctx context.Context, user UserInput) (*AuthOutput, error)
	LogoutUser(ctx context.Context, userId int32, accessToken string, expiresAt time.Time) error
	RefreshTokens(ctx context.Context, refreshToken string) (*AuthOutput, error)
}

type RefreshTokenService interface {
	GetRefreshTokenByUserID(ctx context.Context, qtx repo.Querier, userId int32) (*RefreshTokenOutput, error)
	CreateRefreshToken(ctx context.Context, input CreateRefreshTokenInput) error
	DeleteRefreshTokenByUserID(ctx context.Context, userId int32) error
}

type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (*UserOutput, error)
	GetUserByID(ctx context.Context, id int64) (*UserOutput, error)
	CreateUser(ctx context.Context, qtx repo.Querier, user UserInput) (*UserOutput, error)
}

type GoalService interface {
	ListGoals(ctx context.Context, filter string, userId int32) ([]GoalOutput, error)
	GetGoalByID(ctx context.Context, id int64, userId int32) (*GoalOutput, error)
	CreateGoal(ctx context.Context, qtx repo.Querier, input CreateGoalInput) (*GoalOutput, error)
	UpdateGoal(ctx context.Context, input UpdateGoalInput) (*GoalOutput, error)
	DeleteGoalByID(ctx context.Context, id int64, userId int32) error
}

type RecurringTasksTemplateService interface {
	ListRecurringTasksTemplates(ctx context.Context, userId int32) ([]RecurringTasksTemplateOutput, error)
	GetRecurringTasksTemplateByID(ctx context.Context, id int64, userId int32) (*RecurringTasksTemplateOutput, error)
	CreateRecurringTasksTemplate(ctx context.Context, input CreateRecurringTasksTemplateInput) (*RecurringTasksTemplateOutput, error)
	UpdateRecurringTasksTemplateByID(ctx context.Context, dbTemplate RecurringTasksTemplateOutput, updatingTemplate UpdateRecurringTasksTemplateInput) (*RecurringTasksTemplateOutput, error)
	DeleteRecurringTasksTemplateByID(ctx context.Context, id int64, userId int32) error
	ListRecurringTasksTemplatesDueForGeneration(ctx context.Context, qtx repo.Querier) ([]RecurringTasksTemplateOutput, error)
	UpdateLastGeneratedDateInRecurringTasksTemplateByID(ctx context.Context, qtx repo.Querier, updatingTemplate UpdateLastGeneratedDateInRecurringTasksTemplateInput) error
}

type TaskService interface {
	ListTasksByGoalID(ctx context.Context, userId int32, goalId int32) ([]TaskOutput, error)
	ListInboxTasks(ctx context.Context, userId int32) ([]TaskOutput, error)
	ListTasksByPeriod(ctx context.Context, period GetTasksByPeriodInput) ([]TaskOutput, error)
	ListTasks(ctx context.Context, userId int32) ([]TaskOutput, error)
	GetTaskByID(ctx context.Context, id int64, userId int32) (*TaskOutput, error)
	CreateTask(ctx context.Context, input CreateTaskInput) (*TaskOutput, error)
	UpdateTask(ctx context.Context, dbTask TaskOutput, updatingTask UpdateTaskInput) (*TaskOutput, error)
	AnalyzeForToday(ctx context.Context, userId int32) (*TodayProgressOutput, error)
	DeleteTaskByID(ctx context.Context, id int64, userId int32) error
	DeleteFutureTasksByRecurringTasksTemplateID(ctx context.Context, templateId int64) error
	CreateTasksByRecurringTasksTemplatesDueForGeneration(ctx context.Context, qtx repo.Querier) error
	CreateTasksByRecurringTasksTemplate(ctx context.Context, qtx repo.Querier, template RecurringTasksTemplateOutput) error
}

type AuthCacheRepo interface {
	BlockToken(ctx context.Context, tokenID string, duration time.Duration) error
	IsTokenBlocked(ctx context.Context, tokenID string) (bool, error)
}
