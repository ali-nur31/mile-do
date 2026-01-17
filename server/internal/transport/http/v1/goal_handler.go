package v1

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/ali-nur31/mile-do/internal/domain"
	"github.com/ali-nur31/mile-do/internal/transport/http/v1/dto"
	"github.com/ali-nur31/mile-do/pkg/validator"
	"github.com/labstack/echo/v4"
)

type GoalHandler struct {
	service domain.GoalService
}

func NewGoalHandler(service domain.GoalService) *GoalHandler {
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
// @Success      200  {object}  dto.ListGoalsResponse
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /goals/ [get]
func (h *GoalHandler) GetGoals(c echo.Context) error {
	param := c.QueryParam("type")

	claims, err := GetCurrentClaimsFromCtx(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	goals, err := h.service.ListGoals(c.Request().Context(), param, int32(claims.ID))
	if err != nil {
		slog.Error("failed on getting goals", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	return c.JSON(http.StatusOK, dto.ToListGoalsResponse(goals))
}

// GetGoalByID godoc
// @Summary      get goal by :id
// @Description  get goal by :id
// @Tags         goals
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int64 true "Goal ID"
// @Success      200  {object}  dto.GoalResponse
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      404  {object}  map[string]string "Not Found"
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /goals/{id} [get]
func (h *GoalHandler) GetGoalByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": err.Error()})
	}

	claims, err := GetCurrentClaimsFromCtx(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	goal, err := h.service.GetGoalByID(c.Request().Context(), int64(id), int32(claims.ID))
	if err != nil {
		slog.Error("failed on getting goal by id", "error", err)
		return c.JSON(http.StatusNotFound, map[string]string{"message": "not found", "error": err.Error()})
	}

	return c.JSON(http.StatusOK, dto.ToGoalResponse(goal))
}

// CreateGoal godoc
// @Summary      create new goal
// @Description  create new unique goal
// @Tags         goals
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input body dto.CreateGoalRequest true "Goal Info"
// @Success      201  {object}  dto.GoalResponse
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /goals/ [post]
func (h *GoalHandler) CreateGoal(c echo.Context) error {
	claims, err := GetCurrentClaimsFromCtx(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	var request dto.CreateGoalRequest
	if err = c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": err.Error()})
	}

	if validateErrors := validator.ValidateStruct(request); validateErrors != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "validation failed", "details": validateErrors})
	}

	goal := domain.CreateGoalInput{
		UserID:       int32(claims.ID),
		Title:        request.Title,
		Color:        request.Color,
		CategoryType: request.CategoryType,
	}

	outGoal, err := h.service.CreateGoal(c.Request().Context(), nil, goal)
	if err != nil {
		slog.Error("failed on creating goal", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	return c.JSON(http.StatusCreated, dto.ToGoalResponse(outGoal))
}

// UpdateGoal godoc
// @Summary      update goal
// @Description  update existing goal
// @Tags         goals
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input body dto.UpdateGoalRequest true "New Goal Info"
// @Success      200  {object}  dto.GoalResponse
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /goals/ [patch]
func (h *GoalHandler) UpdateGoal(c echo.Context) error {
	claims, err := GetCurrentClaimsFromCtx(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	var request dto.UpdateGoalRequest
	if err = c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": err.Error()})
	}

	if validateErrors := validator.ValidateStruct(request); validateErrors != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "validation failed", "details": validateErrors})
	}

	outGoal, err := h.service.UpdateGoal(c.Request().Context(), domain.UpdateGoalInput{
		ID:           request.ID,
		UserID:       int32(claims.ID),
		Title:        request.Title,
		Color:        request.Color,
		CategoryType: request.CategoryType,
		IsArchived:   request.IsArchived,
	})
	if err != nil {
		slog.Error("failed on updating goal", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	return c.JSON(http.StatusOK, dto.ToGoalResponse(outGoal))
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
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Failure      400  {object}  map[string]string "Bad Request"
// @Router       /goals/{id} [delete]
func (h *GoalHandler) DeleteGoalByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": err.Error()})
	}

	claims, err := GetCurrentClaimsFromCtx(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	goal, err := h.service.GetGoalByID(c.Request().Context(), int64(id), int32(claims.ID))
	if err != nil {
		slog.Error("failed on getting goal by id for deletion", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	if strings.EqualFold(goal.Title, "routine") || strings.EqualFold(goal.Title, "other") {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": "cannot delete default tasks"})
	}

	err = h.service.DeleteGoalByID(c.Request().Context(), int64(id), int32(claims.ID))
	if err != nil {
		slog.Error("failed on deleting goal by id", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "goal has been removed"})
}
