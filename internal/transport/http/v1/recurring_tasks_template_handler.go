package v1

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/ali-nur31/mile-do/internal/domain"
	"github.com/ali-nur31/mile-do/internal/service"
	"github.com/ali-nur31/mile-do/internal/transport/http/v1/dto"
	"github.com/labstack/echo/v4"
)

type RecurringTasksTemplateHandler struct {
	service service.RecurringTasksTemplateService
}

func NewRecurringTasksTemplateHandler(service service.RecurringTasksTemplateService) *RecurringTasksTemplateHandler {
	return &RecurringTasksTemplateHandler{
		service: service,
	}
}

// GetRecurringTasksTemplates godoc
// @Summary      get recurring tasks templates
// @Description  get recurring tasks templates
// @Tags         recurring-tasks-templates
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dto.ListRecurringTasksTemplatesResponse
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /recurring-tasks-templates/ [get]
func (h *RecurringTasksTemplateHandler) GetRecurringTasksTemplates(c echo.Context) error {
	userId, err := GetCurrentUserIdFromCtx(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	templates, err := h.service.ListRecurringTasksTemplates(c.Request().Context(), userId)
	if err != nil {
		slog.Error("failed on getting recurring tasks templates", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	return c.JSON(http.StatusOK, dto.ToListRecurringTasksTemplatesResponse(templates))
}

// GetRecurringTasksTemplateByID godoc
// @Summary      get recurring tasks template by :id
// @Description  get recurring tasks template by :id
// @Tags         recurring-tasks-templates
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int64 true "Recurring Tasks Template ID"
// @Success      200  {object}  dto.RecurringTasksTemplateResponse
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /recurring-tasks-templates/{id} [get]
func (h *RecurringTasksTemplateHandler) GetRecurringTasksTemplateByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": err.Error()})
	}

	userId, err := GetCurrentUserIdFromCtx(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	template, err := h.service.GetRecurringTasksTemplateByID(c.Request().Context(), int64(id), userId)
	if err != nil {
		slog.Error("failed on getting recurring tasks template by id", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	return c.JSON(http.StatusOK, dto.ToRecurringTasksTemplateResponse(template))
}

// CreateRecurringTasksTemplate godoc
// @Summary      create new recurring tasks template
// @Description  create new recurring tasks template
// @Tags         recurring-tasks-templates
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input body updateRecurringTasksTemplateRequest true "Recurring Tasks Template Info"
// @Success      201  {object}  dto.RecurringTasksTemplateResponse
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /recurring-tasks-templates/ [post]
func (h *RecurringTasksTemplateHandler) CreateRecurringTasksTemplate(c echo.Context) error {
	userId, err := GetCurrentUserIdFromCtx(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	var request dto.UpdateRecurringTasksTemplateRequest
	if err = c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": err.Error()})
	}

	var duration time.Duration
	var startDatetime time.Time
	if request.ScheduledDatetime != "" || (request.ScheduledDatetime != "" && request.ScheduledEndTime != "") {
		startDatetime, duration, err = convertDateTimeAndTime(request.ScheduledDatetime, request.ScheduledEndTime)
		if err != nil {
			slog.Error("failed on creating recurring tasks template", "error", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": err.Error()})
		}
	}

	template := domain.CreateRecurringTasksTemplateInput{
		UserID:            userId,
		GoalID:            request.GoalID,
		Title:             request.Title,
		ScheduledDatetime: startDatetime,
		HasTime:           request.HasTime,
		DurationMinutes:   int32(duration),
		RecurrenceRrule:   request.RecurrenceRrule,
	}

	outTemplate, err := h.service.CreateRecurringTasksTemplate(c.Request().Context(), template)
	if err != nil {
		slog.Error("failed on creating recurring tasks template", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	return c.JSON(http.StatusCreated, dto.ToRecurringTasksTemplateResponse(outTemplate))
}

// UpdateRecurringTasksTemplateByID godoc
// @Summary      update recurring tasks template by :id
// @Description  update existing recurring tasks template by :id
// @Tags         recurring-tasks-templates
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int64 true "Recurring Tasks Template ID"
// @Param        input body dto.UpdateRecurringTasksTemplateRequest true "New Recurring Tasks Template Info"
// @Success      200  {object}  dto.RecurringTasksTemplateResponse
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      404  {object}  map[string]string "Not Found"
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /recurring-tasks-templates/{id} [patch]
func (h *RecurringTasksTemplateHandler) UpdateRecurringTasksTemplateByID(c echo.Context) error {
	userId, err := GetCurrentUserIdFromCtx(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	templateId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": err.Error()})
	}

	var request dto.UpdateRecurringTasksTemplateRequest
	if err = c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": err.Error()})
	}

	dbTemplate, err := h.service.GetRecurringTasksTemplateByID(c.Request().Context(), int64(templateId), userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "cannot find recurring tasks template with provided id", "error": err.Error()})
	}

	startDatetime, duration, err := convertDateTimeAndTime(request.ScheduledDatetime, request.ScheduledEndTime)
	if err != nil {
		slog.Error("failed on updating recurring tasks template by id", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": err.Error()})
	}

	outTemplate, err := h.service.UpdateRecurringTasksTemplateByID(c.Request().Context(), *dbTemplate, domain.UpdateRecurringTasksTemplateInput{
		ID:                int64(templateId),
		UserID:            userId,
		GoalID:            request.GoalID,
		Title:             request.Title,
		ScheduledDatetime: startDatetime,
		HasTime:           request.HasTime,
		DurationMinutes:   int32(duration),
		RecurrenceRrule:   request.RecurrenceRrule,
	})
	if err != nil {
		slog.Error("failed on updating recurring tasks template by id", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	return c.JSON(http.StatusOK, dto.ToRecurringTasksTemplateResponse(outTemplate))
}

// DeleteRecurringTasksTemplateByID godoc
// @Summary      delete recurring tasks template by :id
// @Description  delete recurring tasks template by :id
// @Tags         recurring-tasks-templates
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int64 true "Recurring Tasks Template ID"
// @Success      201  {string}  map[string]string "Recurring tasks template has been removed"
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      404  {object}  map[string]string "Not Found"
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /recurring-tasks-templates/{id} [delete]
func (h *RecurringTasksTemplateHandler) DeleteRecurringTasksTemplateByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": err.Error()})
	}

	userId, err := GetCurrentUserIdFromCtx(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	err = h.service.DeleteRecurringTasksTemplateByID(c.Request().Context(), int64(id), userId)
	if err != nil {
		slog.Error("failed on deleting recurring tasks template by id", "error", err)
		return c.JSON(http.StatusNotFound, map[string]string{"message": "recurring tasks template not found", "error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "recurring tasks template has been removed"})
}

func convertDateTimeAndTime(startDateTimeString, endTimeString string) (time.Time, time.Duration, error) {
	var startTime time.Time
	duration := 15 * time.Minute

	if startDateTimeString == "" {
		return time.Time{}, duration, fmt.Errorf("start date time is empty")
	}

	startDateTime, err := time.Parse(dateTimeLayout, startDateTimeString)
	if err != nil {
		startDateTime, err = time.Parse(dateLayout, startDateTimeString)
		if err != nil {
			return time.Time{}, duration, fmt.Errorf("invalid start datetime format: %v", err)
		}
	}

	startTime, _ = time.Parse(timeLayout, startDateTime.Format(timeLayout))

	if endTimeString != "" {
		var endTime time.Time

		endTime, err = time.Parse(timeLayout, endTimeString)
		if err != nil {
			return time.Time{}, duration, fmt.Errorf("invalid end time format: %v", err)
		}

		if !endTime.IsZero() && endTime.After(startTime) {
			duration = endTime.Sub(startTime)
		}
	}

	return startDateTime, duration, nil
}
