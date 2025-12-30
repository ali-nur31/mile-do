package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ali-nur31/mile-do/internal/domain"
	"github.com/ali-nur31/mile-do/internal/service"
	"github.com/labstack/echo/v4"
)

type taskData struct {
	ID              int64     `json:"id"`
	GoalID          int32     `json:"goal_id"`
	Title           string    `json:"title"`
	IsDone          bool      `json:"is_done"`
	ScheduledDate   string    `json:"scheduled_date"`
	ScheduledTime   string    `json:"scheduled_time"`
	DurationMinutes int       `json:"duration_minutes"`
	RescheduleCount int32     `json:"reschedule_count"`
	CreatedAt       time.Time `json:"created_at"`
}

type listTasksResponse struct {
	UserID   int32      `json:"user_id"`
	TaskData []taskData `json:"task_data"`
}

type createTaskRequest struct {
	UserID               int32  `json:"user_id"`
	GoalID               int32  `json:"goal_id"`
	Title                string `json:"title"`
	ScheduledDateTime    string `json:"scheduled_date_time"`
	ScheduledEndDateTime string `json:"scheduled_end_date_time"`
}

type updateTaskRequest struct {
	ID                   int64  `json:"id"`
	UserID               int32  `json:"user_id"`
	GoalID               int32  `json:"goal_id"`
	Title                string `json:"title"`
	IsDone               bool   `json:"is_done"`
	ScheduledDateTime    string `json:"scheduled_date_time"`
	ScheduledEndDateTime string `json:"scheduled_end_date_time"`
}

type taskResponse struct {
	ID              int64     `json:"id"`
	UserID          int32     `json:"user_id"`
	GoalID          int32     `json:"goal_id"`
	Title           string    `json:"title"`
	IsDone          bool      `json:"is_done"`
	ScheduledDate   string    `json:"scheduled_date"`
	ScheduledTime   string    `json:"scheduled_time"`
	DurationMinutes int       `json:"duration_minutes"`
	RescheduleCount int32     `json:"reschedule_count"`
	CreatedAt       time.Time `json:"created_at"`
}

type TaskHandler struct {
	service service.TaskService
}

func NewTaskHandler(service service.TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

// GetTasks godoc
// @Summary      get tasks
// @Description  get list of tasks
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        type query string false "get goals by type"
// @Success      302  {object}  listTasksResponse
// @Failure      404  {object}  map[string]string "Not Found"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /tasks/ [get]
func (h *TaskHandler) GetTasks(c echo.Context) error {
	userId, err := getCurrentUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}

	tasks, err := h.service.ListTasks(c.Request().Context(), userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	var outTasks listTasksResponse
	outTasks.UserID = userId
	outTasks.TaskData = make([]taskData, len(*tasks))

	for index, task := range *tasks {
		outTasks.TaskData[index].ID = task.ID
		outTasks.TaskData[index].GoalID = task.GoalID
		outTasks.TaskData[index].Title = task.Title
		outTasks.TaskData[index].IsDone = task.IsDone
		outTasks.TaskData[index].ScheduledDate = task.ScheduledDate.Format("2025-31-12")
		outTasks.TaskData[index].ScheduledTime = task.ScheduledTime.Format("15:10")
		outTasks.TaskData[index].DurationMinutes = task.DurationMinutes
		outTasks.TaskData[index].RescheduleCount = task.RescheduleCount
		outTasks.TaskData[index].CreatedAt = task.CreatedAt
	}

	return c.JSON(http.StatusFound, outTasks)
}

// GetTaskByID godoc
// @Summary      get task by :id
// @Description  get task by :id
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int64 true "Goal ID"
// @Success      302  {object}  taskResponse
// @Failure      404  {object}  map[string]string "Not Found"
// @Failure      400  {object}  map[string]string "Bad Request"
// @Router       /tasks/{id} [get]
func (h *TaskHandler) GetTaskByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	userId, err := getCurrentUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}

	task, err := h.service.GetTaskByID(c.Request().Context(), int64(id), userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusFound, taskResponse{
		ID:              task.ID,
		UserID:          task.UserID,
		GoalID:          task.GoalID,
		Title:           task.Title,
		IsDone:          task.IsDone,
		ScheduledDate:   task.ScheduledDate.Format("2025-31-12"),
		ScheduledTime:   task.ScheduledTime.Format("15:10"),
		DurationMinutes: task.DurationMinutes,
		RescheduleCount: task.RescheduleCount,
		CreatedAt:       task.CreatedAt,
	})
}

// CreateTask godoc
// @Summary      create new task
// @Description  create new task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input body createTaskRequest true "Task Info"
// @Success      201  {object}  taskResponse
// @Failure      404  {object}  map[string]string "Not Found"
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /tasks/ [post]
func (h *TaskHandler) CreateTask(c echo.Context) error {
	userId, err := getCurrentUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}

	var request createTaskRequest
	if err = c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": err.Error()})
	}

	scheduledDate, scheduledTime, duration, err := convertDates(request.ScheduledDateTime, request.ScheduledEndDateTime)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": err.Error()})
	}

	task := domain.CreateTaskInput{
		UserID:          userId,
		GoalID:          request.GoalID,
		Title:           request.Title,
		ScheduledDate:   scheduledDate,
		ScheduledTime:   scheduledTime,
		DurationMinutes: int16(duration),
	}

	outTask, err := h.service.CreateTask(c.Request().Context(), task)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, &outTask)
}

// UpdateTask godoc
// @Summary      update task
// @Description  update existing task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input body updateTaskRequest true "New Task Info"
// @Success      200  {object}  taskResponse
// @Failure      404  {object}  map[string]string "Not Found"
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /tasks/ [patch]
func (h *TaskHandler) UpdateTask(c echo.Context) error {
	userId, err := getCurrentUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}

	var request updateTaskRequest
	if err = c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	dbTask, err := h.service.GetTaskByID(c.Request().Context(), request.ID, userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "cannot find task with provided id", "error": err.Error()})
	}

	scheduledDate, scheduledTime, duration, err := convertDates(request.ScheduledDateTime, request.ScheduledEndDateTime)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": err.Error()})
	}

	outTask, err := h.service.UpdateTask(c.Request().Context(), *dbTask, domain.UpdateTask{
		ID:              request.ID,
		UserID:          userId,
		GoalID:          request.GoalID,
		Title:           request.Title,
		IsDone:          request.IsDone,
		ScheduledDate:   scheduledDate,
		ScheduledTime:   scheduledTime,
		DurationMinutes: int16(duration),
		RescheduleCount: dbTask.RescheduleCount,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, &outTask)
}

// DeleteTaskByID godoc
// @Summary      delete task by :id
// @Description  delete task by :id
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int64 true "Task ID"
// @Success      201  {string}  map[string]string "goal has been removed"
// @Failure      404  {object}  map[string]string "Not Found"
// @Failure      400  {object}  map[string]string "Bad Request"
// @Router       /tasks/{id} [delete]
func (h *TaskHandler) DeleteTaskByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	userId, err := getCurrentUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}

	err = h.service.DeleteTaskByID(c.Request().Context(), int64(id), userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "task has been removed"})
}

func convertDates(scheduledDateTimeString, scheduledEndDateTimeString string) (time.Time, time.Time, time.Duration, error) {
	var scheduledDateTime, scheduledDate, scheduledTime time.Time
	duration := 15 * time.Minute

	scheduledDateTime, err := time.Parse("2025-31-12 15:10", scheduledDateTimeString)
	if err != nil {
		return time.Time{}, time.Time{}, duration, err
	}

	if !scheduledDateTime.IsZero() {
		scheduledDate, err = time.Parse("2025-31-12", scheduledDateTimeString)
		if err != nil {
			return time.Time{}, time.Time{}, duration, err
		}

		scheduledTime, err = time.Parse("15:10", scheduledDateTimeString)
		if err != nil {
			return time.Time{}, time.Time{}, duration, err
		}
	}

	scheduledEndDateTime, err := time.Parse("2025-31-12 15:10", scheduledEndDateTimeString)
	if err != nil {
		return time.Time{}, time.Time{}, duration, err
	}

	if !scheduledEndDateTime.IsZero() && scheduledEndDateTime.Hour() != 0 && scheduledEndDateTime.Minute() != 0 {
		duration = scheduledEndDateTime.Sub(scheduledDateTime)
	}

	return scheduledDate, scheduledTime, duration, nil
}
