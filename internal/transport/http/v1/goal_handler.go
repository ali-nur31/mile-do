package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ali-nur31/mile-do/internal/domain"
	"github.com/ali-nur31/mile-do/internal/service"
	"github.com/labstack/echo/v4"
)

type goalData struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Color        string    `json:"color"`
	CategoryType string    `json:"category_type"`
	IsArchived   bool      `json:"is_archived"`
	CreatedAt    time.Time `json:"created_at"`
}

type listGoalsResponse struct {
	UserID int32      `json:"user_id"`
	Data   []goalData `json:"data"`
}

type createGoalRequest struct {
	Title        string `json:"title"`
	Color        string `json:"color"`
	CategoryType string `json:"category_type"`
}

type updateGoalRequest struct {
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	Color        string `json:"color"`
	CategoryType string `json:"category_type"`
	IsArchived   bool   `json:"is_archived"`
}

type updateGoalResponse struct {
	ID           int64  `json:"id"`
	UserID       int32  `json:"user_id"`
	Title        string `json:"title"`
	Color        string `json:"color"`
	CategoryType string `json:"category_type"`
	IsArchived   bool   `json:"is_archived"`
}

type goalResponse struct {
	ID           int64     `json:"id"`
	UserID       int32     `json:"user_id"`
	Title        string    `json:"title"`
	Color        string    `json:"color"`
	CategoryType string    `json:"category_type"`
	IsArchived   bool      `json:"is_archived"`
	CreatedAt    time.Time `json:"created_at"`
}

type GoalHandler struct {
	service service.GoalService
}

func NewGoalHandler(service service.GoalService) *GoalHandler {
	return &GoalHandler{
		service: service,
	}
}

// GetGoals godoc
// @Summary      get goals
// @Description  get list of goals
// @Tags         goals
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        type query string false "get goals by type"
// @Success      302  {object}  listGoalsResponse
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /goals/ [get]
func (h *GoalHandler) GetGoals(c echo.Context) error {
	param := c.QueryParam("type")

	userId, err := getCurrentUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	goals, err := h.service.ListGoals(c.Request().Context(), param, userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	var outGoals listGoalsResponse
	outGoals.UserID = userId
	outGoals.Data = make([]goalData, len(*goals))

	for index, goal := range *goals {
		outGoals.Data[index].ID = goal.ID
		outGoals.Data[index].Title = goal.Title
		outGoals.Data[index].Color = goal.Color
		outGoals.Data[index].CategoryType = goal.CategoryType
		outGoals.Data[index].IsArchived = goal.IsArchived
		outGoals.Data[index].CreatedAt = goal.CreatedAt
	}

	return c.JSON(http.StatusFound, outGoals)
}

// GetGoalByID godoc
// @Summary      get goal by :id
// @Description  get goal by :id
// @Tags         goals
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int64 true "Goal ID"
// @Success      200  {object}  goalResponse
// @Failure      404  {object}  map[string]string "Not Found"
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /goals/{id} [get]
func (h *GoalHandler) GetGoalByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": err.Error()})
	}

	userId, err := getCurrentUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	goal, err := h.service.GetGoalByID(c.Request().Context(), int64(id), userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "not found", "error": err.Error()})
	}

	return c.JSON(http.StatusOK, goalResponse{
		ID:           goal.ID,
		UserID:       userId,
		Title:        goal.Title,
		Color:        goal.Color,
		CategoryType: goal.CategoryType,
		IsArchived:   goal.IsArchived,
		CreatedAt:    goal.CreatedAt,
	})
}

// CreateGoal godoc
// @Summary      create new goal
// @Description  create new unique goal
// @Tags         goals
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input body createGoalRequest true "Goal Info"
// @Success      201  {object}  goalResponse
// @Failure      404  {object}  map[string]string "Not Found"
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /goals/ [post]
func (h *GoalHandler) CreateGoal(c echo.Context) error {
	userId, err := getCurrentUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	var request createGoalRequest
	if err = c.Bind(&request); err != nil {
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
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	return c.JSON(http.StatusCreated, goalResponse{
		ID:           outGoal.ID,
		UserID:       userId,
		Title:        outGoal.Title,
		Color:        outGoal.Color,
		CategoryType: outGoal.CategoryType,
		IsArchived:   outGoal.IsArchived,
		CreatedAt:    outGoal.CreatedAt,
	})
}

// UpdateGoal godoc
// @Summary      update goal
// @Description  update existing goal
// @Tags         goals
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input body updateGoalRequest true "New Goal Info"
// @Success      200  {object}  updateGoalResponse
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /goals/ [patch]
func (h *GoalHandler) UpdateGoal(c echo.Context) error {
	userId, err := getCurrentUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	var request updateGoalRequest
	if err = c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": err.Error()})
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
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	return c.JSON(http.StatusOK, goalResponse{
		ID:           outGoal.ID,
		UserID:       userId,
		Title:        outGoal.Title,
		Color:        outGoal.Color,
		CategoryType: outGoal.CategoryType,
		IsArchived:   outGoal.IsArchived,
	})
}

// DeleteGoalByID godoc
// @Summary      delete goal by :id
// @Description  delete goal by :id
// @Tags         goals
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int64 true "Goal ID"
// @Success      201  {string}  map[string]string "goal has been removed"
// @Failure      404  {object}  map[string]string "Not Found"
// @Failure      400  {object}  map[string]string "Bad Request"
// @Router       /goals/{id} [delete]
func (h *GoalHandler) DeleteGoalByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": err.Error()})
	}

	userId, err := getCurrentUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	err = h.service.DeleteGoalByID(c.Request().Context(), int64(id), userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "goal not found", "error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "goal has been removed"})
}
