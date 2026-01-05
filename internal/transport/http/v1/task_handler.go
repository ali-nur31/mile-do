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

const dateTimeLayout = "2006-01-02 15:04"
const dateLayout = "2006-01-02"
const timeLayout = "15:04"

type TaskHandler struct {
	service service.TaskService
}

func NewTaskHandler(service service.TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

// GetTasksByGoalID godoc
// @Summary      get tasks by :goal_id
// @Description  get tasks by :goal_id
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int64 true "Goal ID"
// @Success      200  {object}  dto.ListTasksResponse
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /goals/{id}/tasks [get]
func (h *TaskHandler) GetTasksByGoalID(c echo.Context) error {
	goalId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": err.Error()})
	}

	userId, err := GetCurrentUserIdFromCtx(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	tasks, err := h.service.ListTasksByGoalID(c.Request().Context(), userId, int32(goalId))
	if err != nil {
		slog.Error("failed on getting tasks by goal id", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	return c.JSON(http.StatusOK, dto.ToListTasksResponse(tasks))
}

// GetInboxTasks godoc
// @Summary      get inbox tasks
// @Description  get tasks without date
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dto.ListTasksResponse
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /tasks/inbox [get]
func (h *TaskHandler) GetInboxTasks(c echo.Context) error {
	userId, err := GetCurrentUserIdFromCtx(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	tasks, err := h.service.ListInboxTasks(c.Request().Context(), userId)
	if err != nil {
		slog.Error("failed on getting inbox tasks", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	return c.JSON(http.StatusOK, dto.ToListTasksResponse(tasks))
}

// GetTasksByPeriod godoc
// @Summary      get tasks by period
// @Description  get tasks by period
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        after_date query string false "tasks after specific date"
// @Param        before_date query string false "tasks before specific date"
// @Success      200  {object}  dto.ListTasksResponse
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /tasks/period [get]
func (h *TaskHandler) GetTasksByPeriod(c echo.Context) error {
	afterDateParam := c.QueryParam("after_date")
	beforeDateParam := c.QueryParam("before_date")
	if afterDateParam == "" || beforeDateParam == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": "at least after_date or before_date must be present"})
	}

	afterDate, err := time.Parse(dateLayout, afterDateParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request, after_date must be in YYYY-DD-MM format", "error": err.Error()})
	}

	beforeDate, err := time.Parse(dateLayout, beforeDateParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request, before_date must be in YYYY-DD-MM format", "error": err.Error()})
	}

	userId, err := GetCurrentUserIdFromCtx(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	tasks, err := h.service.ListTasksByPeriod(c.Request().Context(), domain.GetTasksByPeriodInput{
		UserID:     userId,
		AfterDate:  afterDate,
		BeforeDate: beforeDate,
	})
	if err != nil {
		slog.Error("failed on getting tasks by period", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	return c.JSON(http.StatusOK, dto.ToListTasksResponse(tasks))
}

// GetTasks godoc
// @Summary      get tasks
// @Description  get all tasks
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dto.ListTasksResponse
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /tasks/ [get]
func (h *TaskHandler) GetTasks(c echo.Context) error {
	userId, err := GetCurrentUserIdFromCtx(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	tasks, err := h.service.ListTasks(c.Request().Context(), userId)
	if err != nil {
		slog.Error("failed on getting all tasks", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	return c.JSON(http.StatusOK, dto.ToListTasksResponse(tasks))
}

// GetTaskByID godoc
// @Summary      get task by :id
// @Description  get task by :id
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int64 true "Task ID"
// @Success      200  {object}  dto.TaskResponse
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /tasks/{id} [get]
func (h *TaskHandler) GetTaskByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": err.Error()})
	}

	userId, err := GetCurrentUserIdFromCtx(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	task, err := h.service.GetTaskByID(c.Request().Context(), int64(id), userId)
	if err != nil {
		slog.Error("failed on getting task by id", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	return c.JSON(http.StatusOK, dto.ToTaskResponse(task))
}

// CreateTask godoc
// @Summary      create new task
// @Description  create new task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input body dto.CreateTaskRequest true "Task Info"
// @Success      201  {object}  dto.TaskResponse
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /tasks/ [post]
func (h *TaskHandler) CreateTask(c echo.Context) error {
	userId, err := GetCurrentUserIdFromCtx(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	var request dto.CreateTaskRequest
	if err = c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": err.Error()})
	}

	var scheduledDate, scheduledTime time.Time
	var duration time.Duration
	var hasTime bool
	if request.ScheduledDateTime != "" || (request.ScheduledDateTime != "" && request.ScheduledEndDateTime != "") {
		scheduledDate, scheduledTime, hasTime, duration, err = convertDateTimes(request.ScheduledDateTime, request.ScheduledEndDateTime)
		if err != nil {
			slog.Error("failed on creating task", "error", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": err.Error()})
		}
	}

	task := domain.CreateTaskInput{
		UserID:          userId,
		GoalID:          request.GoalID,
		Title:           request.Title,
		ScheduledDate:   scheduledDate,
		ScheduledTime:   scheduledTime,
		HasTime:         hasTime,
		DurationMinutes: duration,
	}

	outTask, err := h.service.CreateTask(c.Request().Context(), task)
	if err != nil {
		slog.Error("failed on creating task", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	return c.JSON(http.StatusCreated, dto.ToTaskResponse(outTask))
}

// UpdateTask godoc
// @Summary      update task by :id
// @Description  update existing task by :id
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int64 true "Task ID"
// @Param        input body dto.UpdateTaskRequest true "New Task Info"
// @Success      200  {object}  dto.TaskResponse
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      404  {object}  map[string]string "Not Found"
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /tasks/{id} [patch]
func (h *TaskHandler) UpdateTask(c echo.Context) error {
	userId, err := GetCurrentUserIdFromCtx(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	taskId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": err.Error()})
	}

	var request dto.UpdateTaskRequest
	if err = c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": err.Error()})
	}

	dbTask, err := h.service.GetTaskByID(c.Request().Context(), int64(taskId), userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "cannot find task with provided id", "error": err.Error()})
	}

	scheduledDate, scheduledTime, hasTime, duration, err := convertDateTimes(request.ScheduledDateTime, request.ScheduledEndDateTime)
	if err != nil {
		slog.Error("failed on updating task", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": err.Error()})
	}

	outTask, err := h.service.UpdateTask(c.Request().Context(), *dbTask, domain.UpdateTaskInput{
		ID:              int64(taskId),
		UserID:          userId,
		GoalID:          request.GoalID,
		Title:           request.Title,
		IsDone:          request.IsDone,
		ScheduledDate:   scheduledDate,
		ScheduledTime:   scheduledTime,
		HasTime:         hasTime,
		DurationMinutes: duration,
		RescheduleCount: dbTask.RescheduleCount,
	})
	if err != nil {
		slog.Error("failed on updating task", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	return c.JSON(http.StatusOK, dto.ToTaskResponse(outTask))
}

// AnalyzeForToday godoc
// @Summary      get stats for today
// @Description  get count of completed tasks over total tasks for today
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      201  {string}  dto.CountCompletedTasksForTodayResponse
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      400  {object}  map[string]string "Bad Request"
// @Router       /tasks/analyze [get]
func (h *TaskHandler) AnalyzeForToday(c echo.Context) error {
	userId, err := GetCurrentUserIdFromCtx(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	stats, err := h.service.AnalyzeForToday(c.Request().Context(), userId)
	if err != nil {
		slog.Error("failed on getting tasks analysis for today", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	return c.JSON(http.StatusOK, dto.ToCountCompletedTasksForTodayResponse(stats))
}

// DeleteTaskByID godoc
// @Summary      delete task by :id
// @Description  delete task by :id
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int64 true "Task ID"
// @Success      201  {string}  map[string]string "Task has been removed"
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      404  {object}  map[string]string "Not Found"
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /tasks/{id} [delete]
func (h *TaskHandler) DeleteTaskByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": err.Error()})
	}

	userId, err := GetCurrentUserIdFromCtx(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	err = h.service.DeleteTaskByID(c.Request().Context(), int64(id), userId)
	if err != nil {
		slog.Error("failed on deleting task by id", "error", err)
		return c.JSON(http.StatusNotFound, map[string]string{"message": "task not found", "error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "task has been removed"})
}

func convertDateTimes(startDateTimeString, endDateTimeString string) (time.Time, time.Time, bool, time.Duration, error) {
	var startDate, startTime time.Time
	duration := 15 * time.Minute
	hasTime := false

	if startDateTimeString == "" {
		return time.Time{}, time.Time{}, false, duration, fmt.Errorf("start date time is empty")
	}

	startDateTime, err := time.Parse(dateTimeLayout, startDateTimeString)
	if err == nil {
		hasTime = true
	} else {
		startDateTime, err = time.Parse(dateLayout, startDateTimeString)
		if err != nil {
			return time.Time{}, time.Time{}, false, duration, fmt.Errorf("invalid start date time format: %v", err)
		}
	}

	startDate, _ = time.Parse(dateLayout, startDateTime.Format(dateLayout))
	startTime, _ = time.Parse(timeLayout, startDateTime.Format(timeLayout))

	if endDateTimeString != "" && hasTime {
		var endDateTime time.Time

		endDateTime, err = time.Parse(dateTimeLayout, endDateTimeString)
		if err != nil {
			endDateTime, err = time.Parse(dateLayout, endDateTimeString)
			if err != nil {
				return time.Time{}, time.Time{}, false, duration, fmt.Errorf("invalid start date time format: %v", err)
			}
		}

		if !endDateTime.IsZero() && endDateTime.After(startDateTime) {
			duration = endDateTime.Sub(startDateTime)
		}
	}

	return startDate, startTime, hasTime, duration, nil
}
