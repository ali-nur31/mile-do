package dto

import (
	"github.com/ali-nur31/mile-do/internal/domain"
)

type CreateTaskRequest struct {
	UserID               int32  `json:"user_id" validate:"required,gte=0"`
	GoalID               int32  `json:"goal_id" validate:"required,gte=0"`
	Title                string `json:"title" validate:"required,min=3,max=256"`
	ScheduledDateTime    string `json:"scheduled_date_time" validate:"required,min=10"`
	ScheduledEndDateTime string `json:"scheduled_end_date_time" validate:"omitempty,min=10"`
}

type UpdateTaskRequest struct {
	GoalID               int32  `json:"goal_id" validate:"required,gte=0"`
	Title                string `json:"title" validate:"required,min=3,max=256"`
	IsDone               bool   `json:"is_done" validate:"required,oneof=true false"`
	ScheduledDateTime    string `json:"scheduled_date_time" validate:"required,min=10"`
	ScheduledEndDateTime string `json:"scheduled_end_date_time" validate:"omitempty,min=10"`
}

type CountCompletedTasksForTodayResponse struct {
	TotalTasks int32 `json:"total_tasks"`
	Completed  int32 `json:"completed"`
}

func ToCountCompletedTasksForTodayResponse(output *domain.TodayProgressOutput) CountCompletedTasksForTodayResponse {
	return CountCompletedTasksForTodayResponse{
		TotalTasks: output.TotalTasks,
		Completed:  output.CompletedToday,
	}
}

type TaskResponse struct {
	ID              int64  `json:"id"`
	UserID          int32  `json:"user_id"`
	GoalID          int32  `json:"goal_id"`
	Title           string `json:"title"`
	IsDone          bool   `json:"is_done"`
	ScheduledDate   string `json:"scheduled_date"`
	HasTime         bool   `json:"has_time"`
	ScheduledTime   string `json:"scheduled_time"`
	DurationMinutes int32  `json:"duration_minutes"`
	RescheduleCount int32  `json:"reschedule_count"`
	CreatedAt       string `json:"created_at"`
}

func ToTaskResponse(task *domain.TaskOutput) TaskResponse {
	return TaskResponse{
		ID:              task.ID,
		UserID:          task.UserID,
		GoalID:          task.GoalID,
		Title:           task.Title,
		IsDone:          task.IsDone,
		ScheduledDate:   task.ScheduledDate.String(),
		HasTime:         task.HasTime,
		ScheduledTime:   task.ScheduledTime.String(),
		DurationMinutes: task.DurationMinutes,
		RescheduleCount: task.RescheduleCount,
		CreatedAt:       task.CreatedAt.String(),
	}
}

type TaskData struct {
	ID              int64  `json:"id"`
	GoalID          int32  `json:"goal_id"`
	Title           string `json:"title"`
	IsDone          bool   `json:"is_done"`
	ScheduledDate   string `json:"scheduled_date"`
	HasTime         bool   `json:"has_time"`
	ScheduledTime   string `json:"scheduled_time"`
	DurationMinutes int32  `json:"duration_minutes"`
	RescheduleCount int32  `json:"reschedule_count"`
	CreatedAt       string `json:"created_at"`
}

type ListTasksResponse struct {
	UserID   int32      `json:"user_id"`
	TaskData []TaskData `json:"task_data"`
}

func ToListTasksResponse(tasks []domain.TaskOutput) ListTasksResponse {
	taskData := make([]TaskData, len(tasks))

	for index, task := range tasks {
		taskData[index] = TaskData{
			ID:              task.ID,
			GoalID:          task.GoalID,
			Title:           task.Title,
			IsDone:          task.IsDone,
			ScheduledDate:   task.ScheduledDate.String(),
			HasTime:         task.HasTime,
			ScheduledTime:   task.ScheduledTime.String(),
			DurationMinutes: task.DurationMinutes,
			RescheduleCount: task.RescheduleCount,
			CreatedAt:       task.CreatedAt.String(),
		}
	}

	return ListTasksResponse{
		UserID:   tasks[0].UserID,
		TaskData: taskData,
	}
}
