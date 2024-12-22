-- name: InsertUpvote :one
INSERT INTO upvotes (user_id, post_id)
VALUES ($1, $2)
RETURNING user_id, post_id, created_at;

-- name: DeleteUpvote :one
DELETE FROM upvotes
WHERE user_id = $1 AND post_id = $2
RETURNING user_id, post_id, created_at;
