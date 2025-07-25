-- name: CreatePost :one
INSERT INTO posts (user_id, group_id, content)
VALUES ($1, $2, $3)
RETURNING id, user_id, group_id;

-- name: CreateComment :one
INSERT INTO posts (user_id, group_id, content, parent_post_id)
VALUES ($1, $2, $3, $4)
RETURNING id, parent_post_id;

-- name: GetPostByID :one
SELECT * FROM posts WHERE id =  $1;

-- name: GetPostsByGroupID :many
SELECT * FROM posts
WHERE group_id = $1
ORDER BY created_at DESC;

-- name: UpdatePost :one
UPDATE posts
SET content = $1, updated_at = CURRENT_TIMESTAMP
WHERE id = $2
AND is_deleted = FALSE
RETURNING *;

-- name: DeletePost :exec
UPDATE posts
SET is_deleted = TRUE, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: RestorePost :one
UPDATE posts
SET is_deleted = FALSE, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: GetPostsForFeed :many
SELECT * FROM posts
WHERE group_id = $1
AND is_deleted = FALSE
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ResetPosts :exec
TRUNCATE TABLE posts CASCADE;

-- name: DeleteCommentsFromPost :exec
UPDATE posts
SET updated_at = CURRENT_TIMESTAMP, is_deleted = TRUE
WHERE parent_post_id = $1;