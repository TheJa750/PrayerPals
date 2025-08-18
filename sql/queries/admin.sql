-- name: GetKickBanStatus :one
SELECT is_banned, is_kicked, kicked_until
FROM users_groups
WHERE user_id = $1 AND group_id = $2;

-- name: BanUser :exec
UPDATE users_groups
SET is_banned = TRUE, kicked_until = NULL, modded_reason = $3,
    modded_at = NOW(), modded_by = $4
WHERE user_id = $1 AND group_id = $2;

-- name: UnbanUser :exec
UPDATE users_groups
SET is_banned = FALSE, kicked_until = NULL, modded_reason = '',
    modded_at = NOW(), modded_by = $3
WHERE user_id = $1 AND group_id = $2;

-- name: KickUser :exec
UPDATE users_groups
SET is_kicked = TRUE, kicked_until = NOW() + INTERVAL '7 days', modded_reason = $3,
    modded_at = NOW(), modded_by = $4
WHERE user_id = $1 AND group_id = $2;

-- name: RemovePostsByUser :exec
UPDATE posts
SET is_deleted = TRUE, updated_at = NOW()
WHERE user_id = $1 AND group_id = $2 AND is_deleted = FALSE;

-- name: UpdateGroupInviteCode :exec
UPDATE groups
SET invite_code = $2
WHERE id = $1;

-- name: UpdateGroupDescription :exec
UPDATE groups
SET description = $2
WHERE id = $1;

-- name: UpdateGroupRules :exec
UPDATE groups
SET rules_info = $2
WHERE id = $1;