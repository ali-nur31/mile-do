-- name: GetGoalByID :one
SELECT * FROM goals
WHERE id = $1 AND user_id = $2 LIMIT 1;

-- name: ListGoalsByIsArchived :many
SELECT * FROM goals
WHERE is_archived = $1 AND user_id = $2
ORDER BY id;

-- name: ListGoals :many
SELECT * FROM goals
WHERE user_id = $1
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
SET title = $3, color = $4, category_type = $5, is_archived = $6
WHERE id = $1 AND user_id = $2;

-- name: DeleteGoalByID :exec
DELETE FROM goals
WHERE id = $1 AND user_id = $2;