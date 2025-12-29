package v1

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/ali-nur31/mile-do/internal/domain"
	"github.com/ali-nur31/mile-do/internal/service"
	"github.com/labstack/echo/v4"
)

type listGoalsOutput struct {
	UserID int32 `json:"user_id"`
	Data   struct {
		ID           int64     `json:"id"`
		Title        string    `json:"title"`
		Color        string    `json:"color"`
		CategoryType string    `json:"category_type"`
		IsArchived   bool      `json:"is_archived"`
		CreatedAt    time.Time `json:"created_at"`
	} `json:"data"`
}

type createGoalInput struct {
	Title        string `json:"title"`
	Color        string `json:"color"`
	CategoryType string `json:"category_type"`
}

type updateGoalInput struct {
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	Color        string `json:"color"`
	CategoryType string `json:"category_type"`
	IsArchived   bool   `json:"is_archived"`
}

type GoalHandler struct {
	service service.GoalService
}

func NewGoalHandler(service service.GoalService) *GoalHandler {
	return &GoalHandler{
		service: service,
	}
}

func (h *GoalHandler) GetGoals(c echo.Context) error {
	param := c.QueryParam("type")

	userIdFromCtx := c.Get("userId")
	userId, ok := userIdFromCtx.(int32)
	if !ok {
		slog.Error("email in context is not a string", "value", userIdFromCtx)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}

	goals, err := h.service.ListGoals(c.Request().Context(), param)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	var outGoals listGoalsOutput
	outGoals.UserID = userId

	for _, goal := range *goals {
		outGoals.Data.ID = goal.ID
		outGoals.Data.Title = goal.Title
		outGoals.Data.Color = goal.Color
		outGoals.Data.CategoryType = goal.CategoryType
		outGoals.Data.IsArchived = goal.IsArchived
		outGoals.Data.CreatedAt = goal.CreatedAt
	}

	return c.JSON(http.StatusFound, outGoals)
}

func (h *GoalHandler) GetGoalByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	goal, err := h.service.GetGoalByID(c.Request().Context(), int64(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusFound, goal)
}

func (h *GoalHandler) CreateGoal(c echo.Context) error {
	userIdFromCtx := c.Get("userId")
	userId, ok := userIdFromCtx.(int32)
	if !ok {
		slog.Error("email in context is not a string", "value", userIdFromCtx)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}

	var request createGoalInput
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	goal := domain.CreateGoalInput{
		UserID:       userId,
		Title:        request.Title,
		Color:        request.Color,
		CategoryType: request.CategoryType,
	}

	outGoal, err := h.service.CreateGoal(c.Request().Context(), goal)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, &outGoal)
}

func (h *GoalHandler) UpdateGoal(c echo.Context) error {
	userIdFromCtx := c.Get("userId")
	userId, ok := userIdFromCtx.(int32)
	if !ok {
		slog.Error("userId in context is not an int", "value", userIdFromCtx)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}

	var request updateGoalInput
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	outGoal, err := h.service.UpdateGoal(c.Request().Context(), domain.UpdateGoalInput{
		ID:           request.ID,
		UserID:       userId,
		Title:        request.Title,
		Color:        request.Color,
		CategoryType: request.CategoryType,
		IsArchived:   request.IsArchived,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, &outGoal)
}

func (h *GoalHandler) DeleteGoalByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err = h.service.DeleteGoalByID(c.Request().Context(), int64(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, "goal has been removed")
}
