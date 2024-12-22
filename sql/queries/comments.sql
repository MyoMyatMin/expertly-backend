-- name: CreateComment :one
INSERT INTO comments (id, post_id, user_id, parent_comment_id, content, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
RETURNING id, post_id, user_id, parent_comment_id, content, created_at, updated_at;

-- name: GetCommentsByPost :many
SELECT id, post_id, user_id, parent_comment_id, content, created_at, updated_at
FROM comments
WHERE post_id = $1
ORDER BY created_at ASC;

-- name: UpdateComment :one
UPDATE comments SET content = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, post_id, user_id, parent_comment_id, content, created_at, updated_at;

-- name: DeleteComment :exec
DELETE FROM comments
WHERE id = $1;