-- name: CreateGroup :one
INSERT INTO groups (name, description, owner_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetGroupByID :one
SELECT *
FROM groups
WHERE id = $1;

-- name: AddUserToGroup :exec
INSERT INTO users_groups (user_id, group_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: RemoveUserFromGroup :exec
DELETE FROM users_groups
WHERE user_id = $1 AND group_id = $2;