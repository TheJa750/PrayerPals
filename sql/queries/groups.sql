-- name: CreateGroup :one
INSERT INTO groups (name, description, owner_id, invite_code)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetGroupByID :one
SELECT *
FROM groups
WHERE id = $1;

-- name: AddUserToGroup :exec
INSERT INTO users_groups (user_id, group_id, role)
VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING;

-- name: RemoveUserFromGroup :exec
DELETE FROM users_groups
WHERE user_id = $1 AND group_id = $2;

-- name: DeleteGroup :exec
DELETE FROM groups
WHERE id = $1;

-- name: ResetGroups :exec
TRUNCATE TABLE groups CASCADE;

-- name: GetGroupMembersIDs :many
SELECT user_id
FROM users_groups
WHERE group_id = $1;

-- name: GetGroupsForUser :many
SELECT group_id
FROM users_groups
WHERE user_id = $1;

-- name: GetGroupSpecialRoles :many
SELECT user_id, role
FROM users_groups
WHERE group_id = $1 AND role != 'member';

-- name: GetGroupByInviteCode :one
SELECT * FROM groups
WHERE invite_code = $1;

-- name: GetActiveMembers :many
SELECT
    users.id,
    users.username,
    users.email,
    users_groups.role
FROM users_groups
JOIN users ON users.id = users_groups.user_id
WHERE users_groups.group_id = $1
AND NOT users_groups.is_banned
AND NOT users_groups.is_kicked;