-- name: GetRefreshTokenByUserID :one
SELECT * FROM refresh_tokens
WHERE user_id = $1 LIMIT 1;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
    user_id, token_hash, expires_at
) VALUES (
 $1, $2, $3
)
RETURNING *;

-- name: DeleteRefreshTokenByUserID :exec
DELETE FROM refresh_tokens
WHERE user_id = $1;