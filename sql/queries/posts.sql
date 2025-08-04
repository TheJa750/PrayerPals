-- name: CreatePost :one
INSERT INTO posts (user_id, group_id, content)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateComment :one
INSERT INTO posts (user_id, group_id, content, parent_post_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

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
SELECT
    posts.id,
    posts.content,
    posts.user_id,
    posts.group_id,
    posts.created_at,
    users.username,
    COUNT(comments.id) AS comment_count
FROM posts
LEFT JOIN users ON posts.user_id = users.id
LEFT JOIN posts AS comments ON posts.id = comments.parent_post_id
WHERE posts.group_id = $1
AND posts.parent_post_id IS NULL
AND posts.is_deleted = FALSE
GROUP BY posts.id, posts.content, posts.user_id, posts.group_id, posts.created_at, users.username
ORDER BY posts.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ResetPosts :exec
TRUNCATE TABLE posts CASCADE;

-- name: DeleteCommentsFromPost :exec
UPDATE posts
SET updated_at = CURRENT_TIMESTAMP, is_deleted = TRUE
WHERE parent_post_id = $1;

-- name: GetCommentsByPostID :many
SELECT * FROM posts
WHERE parent_post_id = $1
AND is_deleted = FALSE
ORDER BY created_at DESC;