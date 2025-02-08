-- name: CreateSavedPost :one
INSERT INTO saved_posts(user_id, post_id) VALUES ($1, $2) RETURNING *;

-- name: DeleteSavedPost :exec
DELETE FROM saved_posts WHERE user_id = $1 AND post_id = $2;

-- name: ListSavedPostsByID :many
SELECT 
    p.post_id as post_id,
    p.user_id as user_id,
    p.slug as slug,
    p.title as title,
    p.content as content,
    p.created_at as created_at,
    p.updated_at as updated_at,
    COALESCE(u.upvote_count, 0) as upvote_count,
    COALESCE(c.comment_count, 0) as comment_count
FROM posts p
JOIN saved_posts s ON p.post_id = s.post_id
LEFT JOIN (
    SELECT post_id, COUNT(*) as upvote_count
    FROM upvotes
    GROUP BY post_id
) u ON p.post_id = u.post_id
LEFT JOIN (
    SELECT post_id, COUNT(*) as comment_count
    FROM comments
    GROUP BY post_id
) c ON p.post_id = c.post_id
WHERE s.user_id = $1
ORDER BY p.created_at DESC;

