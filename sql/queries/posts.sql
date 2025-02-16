-- name: CreatePost :one
INSERT INTO posts (
    post_id,
    user_id,
    slug,
    title,
    content
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING post_id, user_id, slug, title, content, created_at, updated_at;

-- name: GetPost :one
SELECT 
    post_id, 
    user_id, 
    slug, 
    title, 
    content, 
    created_at, 
    updated_at
FROM posts
WHERE post_id = $1;

-- name: UpdatePost :one
UPDATE posts
SET
    title = $2,
    slug = $3,
    content = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE post_id = $1
RETURNING post_id, user_id, slug, title, content, created_at, updated_at;

-- name: DeletePost :exec
DELETE FROM posts
WHERE post_id = $1;

-- name: ListPosts :many
SELECT 
    p.post_id, 
    p.slug, 
    p.title, 
    p.user_id, 
    p.content, 
    p.created_at, 
    p.updated_at, 
    COALESCE(upvote_counts.count, 0) AS upvote_count, 
    COALESCE(comment_counts.count, 0) AS comment_count
FROM posts p
LEFT JOIN (
    SELECT 
        post_id, 
        COUNT(*) AS count 
    FROM comments 
    GROUP BY post_id
) comment_counts ON p.post_id = comment_counts.post_id
LEFT JOIN (
    SELECT 
        post_id, 
        COUNT(*) AS count 
    FROM upvotes 
    GROUP BY post_id
) upvote_counts ON p.post_id = upvote_counts.post_id
WHERE p.created_at >= NOW() - INTERVAL '30 days'  -- Get posts from last 30 days
ORDER BY 
    (COALESCE(upvote_counts.count, 0) * 2 + COALESCE(comment_counts.count, 0)) DESC, -- Weight upvotes & comments
    p.created_at DESC -- Prioritize newer posts if engagement is similar
LIMIT 20;

-- name: GetPostBySlug :one
SELECT 
    post_id, 
    user_id, 
    slug, 
    title, 
    content, 
    created_at, 
    updated_at
FROM posts
WHERE slug = $1;

-- name: GetPostDetailsByID :one
SELECT 
    p.post_id, 
    p.user_id, 
    p.slug, 
    p.title, 
    p.content, 
    p.created_at, 
    p.updated_at,
    u.name AS author_name,
    u.username AS author_username,
    COALESCE(upvote_counts.count, 0) AS upvote_count,
    COALESCE(comment_counts.count, 0) AS comment_count
FROM posts p
JOIN users u ON p.user_id = u.user_id
LEFT JOIN (
    SELECT 
        post_id, 
        COUNT(*) AS count
    FROM upvotes
    GROUP BY post_id
) upvote_counts ON p.post_id = upvote_counts.post_id
LEFT JOIN (
    SELECT 
        post_id, 
        COUNT(*) AS count
    FROM comments
    GROUP BY post_id
) comment_counts ON p.post_id = comment_counts.post_id
WHERE p.post_id = $1;

-- name: GetPostDetailsForUsersByID :one
SELECT 
    p.post_id, 
    p.user_id, 
    p.slug, 
    p.title, 
    p.content, 
    p.created_at, 
    p.updated_at,
    u.name AS author_name,
    u.username AS author_username,
    COALESCE(upvote_counts.count, 0) AS upvote_count,
    COALESCE(comment_counts.count, 0) AS comment_count,
    EXISTS (
        SELECT 1 
        FROM upvotes 
        WHERE upvotes.post_id = p.post_id  -- Prefix `post_id` with table alias
        AND upvotes.user_id = $2          -- Prefix `user_id` with table alias
    ) AS has_upvoted,
    EXISTS (
        SELECT 1 
        FROM saved_posts 
        WHERE saved_posts.post_id = p.post_id  -- Prefix `post_id` with table alias
        AND saved_posts.user_id = $2            -- Prefix `user_id` with table alias
    ) AS has_saved
FROM posts p
JOIN users u ON p.user_id = u.user_id
LEFT JOIN (
    SELECT 
        post_id, 
        COUNT(*) AS count
    FROM upvotes
    GROUP BY post_id
) upvote_counts ON p.post_id = upvote_counts.post_id
LEFT JOIN (
    SELECT 
        post_id, 
        COUNT(*) AS count
    FROM comments
    GROUP BY post_id
) comment_counts ON p.post_id = comment_counts.post_id
WHERE p.post_id = $1;


-- name: DeletePostBySlug :exec
DELETE FROM posts
WHERE slug = $1;

-- name: PostSearchByKeyword :many
SELECT 
    p.post_id, 
    p.slug, 
    p.title, 
    p.user_id, 
    p.content, 
    p.created_at, 
    p.updated_at, 
    u.name AS author_name,
    u.username AS author_username
FROM posts p
JOIN users u ON p.user_id = u.user_id
WHERE p.title ILIKE '%' || $1 || '%'
ORDER BY 
    p.created_at DESC -- Prioritize newer posts
LIMIT 20;
