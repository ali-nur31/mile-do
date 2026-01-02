-- name: GetRecurringTasksTemplateByID :one
SELECT * FROM recurring_tasks_templates
WHERE id = $1 AND user_id = $2 LIMIT 1;

-- name: ListRecurringTasksTemplates :many
SELECT * FROM recurring_tasks_templates
WHERE user_id = $1
ORDER BY id;

-- name: CreateRecurringTasksTemplate :one
INSERT INTO recurring_tasks_templates (
    user_id, goal_id, title, scheduled_datetime, has_time, duration_minutes, recurrence_rrule
) VALUES (
             $1, $2, $3, $4, $5, $6, $7
         )
    RETURNING *;

-- name: UpdateRecurringTasksTemplateByID :exec
UPDATE recurring_tasks_templates
SET
    goal_id = $3,
    title = $4,
    scheduled_datetime = $5,
    has_time = $6,
    duration_minutes = $7,
    recurrence_rrule = $8,
    last_generated_date = $9
WHERE id = $1 AND user_id = $2;

-- name: DeleteRecurringTasksTemplateByID :exec
DELETE FROM recurring_tasks_templates
WHERE id = $1 AND user_id = $2;

-- name: ListRecurringTasksTemplatesDueForGeneration :many
SELECT * FROM recurring_tasks_templates
WHERE last_generated_date < (current_date + interval '1 month');

-- name: UpdateLastGeneratedDateInRecurringTasksTemplateByID :exec
UPDATE recurring_tasks_templates
SET last_generated_date = $2
WHERE id = $1;