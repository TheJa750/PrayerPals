-- name: CreateUser :one
INSERT INTO users (username, email, hashed_password)
VALUES ($1, $2, $3)
RETURNING id, username, email, created_at, updated_at, is_active;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetUserIDByEmail :one
SELECT id
FROM users
WHERE email = $1;

-- name: GetUserGroups :many
SELECT g.id, g.name, g.description, g.created_at, g.updated_at
FROM groups g
JOIN users_groups ug ON g.id = ug.group_id
WHERE ug.user_id = $1;