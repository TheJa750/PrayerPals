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