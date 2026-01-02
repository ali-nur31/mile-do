-- name: GetTaskByID :one
SELECT * FROM tasks
WHERE id = $1 AND user_id = $2 LIMIT 1;

-- name: ListTasksByGoalID :many
SELECT * FROM tasks
WHERE goal_id = $1 AND user_id = $2
ORDER BY is_done ASC, id DESC;

-- name: ListInboxTasks :many
SELECT * FROM tasks
WHERE scheduled_date IS null AND has_time = false AND is_done = false AND user_id = $1
ORDER BY id DESC;

-- name: ListTasksByDateRange :many
SELECT * FROM tasks
WHERE user_id = $3 AND scheduled_date >= $1 AND scheduled_date <= $2
ORDER BY scheduled_time ASC, id;

-- name: CountCompletedTasksForToday :one
SELECT
    count(*) FILTER (WHERE scheduled_date = current_date)::int AS total_today,
    count(*) FILTER (WHERE scheduled_date = current_date AND is_done = true)::int AS completed_today
FROM tasks
WHERE user_id = $1;

-- name: ListTasks :many
SELECT * FROM tasks
WHERE user_id = $1
ORDER BY id;

-- name: CreateTask :one
INSERT INTO tasks (
    user_id, goal_id, recurring_template_id, title, scheduled_date, has_time, scheduled_time, duration_minutes
) VALUES (
             $1, $2, $3, $4, $5, $6, $7, $8
         )
    RETURNING *;

-- name: UpdateTaskByID :exec
UPDATE tasks
SET
    goal_id = $3,
    recurring_template_id = $4,
    title = $5,
    is_done = $6,
    scheduled_date = $7,
    has_time = $8,
    scheduled_time = $9,
    duration_minutes = $10,
    reschedule_count = $11
WHERE id = $1 AND user_id = $2;

-- name: DeleteTaskByID :exec
DELETE FROM tasks
WHERE id = $1 AND user_id = $2;

-- name: DeleteFutureTasksByRecurringTasksTemplateID :exec
DELETE FROM tasks
WHERE recurring_template_id = $1 AND scheduled_date > current_date AND is_done = false;