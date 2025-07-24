-- name: CreateUserToken :one
INSERT INTO refresh_tokens (token, user_id, expires_at)
VALUES ($1, $2, NOW() + INTERVAL '30 days')
RETURNING *;

-- name: GetUserByToken :one
SELECT * FROM refresh_tokens
WHERE token = $1
AND revoked_at IS NULL;

-- name: RevokeUserToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW()
WHERE token = $1;