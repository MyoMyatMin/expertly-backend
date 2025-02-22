-- name: CreateContributor :one
INSERT INTO Contributors (
    user_id,
    expertise_fields
) VALUES (
    $1,
    $2
)
RETURNING user_id, expertise_fields, created_at;

-- name: GetContributorByUserId :one
SELECT 
    user_id,
    expertise_fields,
    created_at
FROM Contributors
WHERE user_id = $1;

-- name: CheckIfUserIsContributor :one
SELECT EXISTS (
    SELECT 1
    FROM Contributors
    WHERE user_id = $1
);


-- name: GetPostsByContributor :many
SELECT 
    p.post_id, 
    p.user_id, 
    p.slug, 
    p.title, 
    p.content, 
    p.created_at, 
    p.updated_at,
    COUNT(DISTINCT u.user_id) AS upvote_count,
    COUNT(DISTINCT c.comment_id) AS comment_count
FROM posts p
LEFT JOIN upvotes u ON p.post_id = u.post_id
LEFT JOIN comments c ON p.post_id = c.post_id
WHERE p.user_id = $1
GROUP BY p.post_id
ORDER BY p.created_at DESC;
