-- name: CreateUser :one
INSERT INTO users (username, email, hashed_password)
VALUES ($1, $2, $3)
RETURNING id, username, email;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetUserIDByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: GetUserGroupIDs :many
SELECT group_id
FROM users_groups
WHERE user_id = $1;

-- name: ResetUsers :exec
TRUNCATE TABLE users CASCADE;