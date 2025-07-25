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

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: UpdateUser :exec
UPDATE users
SET username = $2, email = $3, hashed_password = $4
WHERE id = $1;

-- name: GetUserGroupRole :one
SELECT role
FROM users_groups
WHERE user_id = $1 AND group_id = $2;

-- name: AdjustUserGroupRole :exec
UPDATE users_groups
SET role = $1
WHERE user_id = $2 AND group_id = $3;