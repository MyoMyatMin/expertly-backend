-- name: CreateSavedPost :one
INSERT INTO saved_posts(user_id, post_id) VALUES ($1, $2) RETURNING *;

-- name: DeleteSavedPost :exec
DELETE FROM saved_posts WHERE user_id = $1 AND post_id = $2;