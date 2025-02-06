-- name: CreateComment :one
INSERT INTO comments (
    comment_id,
    post_id,
    user_id,
    parent_comment_id,
    content
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING comment_id, post_id, user_id, parent_comment_id, content, created_at, updated_at;

-- name: GetCommentsByPost :many
SELECT 
    comment_id, 
    post_id, 
    user_id, 
    parent_comment_id, 
    content, 
    created_at, 
    updated_at
FROM comments
WHERE post_id = $1
ORDER BY created_at ASC;

-- name: UpdateComment :one
UPDATE comments 
SET 
    content = $2, 
    updated_at = CURRENT_TIMESTAMP
WHERE comment_id = $1
RETURNING comment_id, post_id, user_id, parent_comment_id, content, created_at, updated_at;

-- name: DeleteComment :exec
DELETE FROM comments
WHERE comment_id = $1;

-- name: GetCommentByID :one
SELECT 
    comment_id, 
    post_id, 
    user_id, 
    parent_comment_id, 
    content, 
    created_at, 
    updated_at
FROM comments
WHERE comment_id = $1;

