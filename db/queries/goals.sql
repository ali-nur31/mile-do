-- name: GetGoalByID :one
SELECT * FROM goals
WHERE id = $1 LIMIT 1;

-- name: ListGoalsByIsArchived :many
SELECT * FROM goals
WHERE is_archived = $1
ORDER BY id;

-- name: ListGoals :many
SELECT * FROM goals
ORDER BY id;

-- name: CreateGoal :one
INSERT INTO goals (
    user_id, title, color, category_type
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateGoalByID :exec
UPDATE goals
SET title = $2, color = $3, category_type = $4, is_archived = $5
WHERE id = $1;

-- name: DeleteGoalByID :exec
DELETE FROM goals
WHERE id = $1;