-- name: GetTaskByID :one
SELECT * FROM tasks
WHERE id = $1 AND user_id = $2 LIMIT 1;

-- name: ListTasksByGoalID :many
SELECT * FROM tasks
WHERE goal_id = $1 AND user_id = $2
ORDER BY is_done ASC, id DESC;

-- name: ListInboxTasks :many
SELECT * FROM tasks
WHERE scheduled_date IS null AND is_done = false AND user_id = $1
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
    user_id, goal_id, title, scheduled_date, scheduled_time, duration_minutes
) VALUES (
             $1, $2, $3, $4, $5, $6
         )
    RETURNING *;

-- name: UpdateTaskByID :exec
UPDATE tasks
SET
    goal_id = $3,
    title = $4,
    is_done = $5,
    scheduled_date = $6,
    scheduled_time = $7,
    duration_minutes = $8,
    reschedule_count = $9
WHERE id = $1 AND user_id = $2;

-- name: DeleteTaskByID :exec
DELETE FROM tasks
WHERE id = $1 AND user_id = $2;