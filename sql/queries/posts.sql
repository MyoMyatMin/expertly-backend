-- name: CreatePost :one
INSERT INTO posts (id, user_id, slug, title, content)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING id, user_id, slug, title, content, created_at, updated_at;

-- name: GetPost :one
SELECT id, user_id, slug, title, content, created_at, updated_at
FROM posts
WHERE id = $1;

-- name: UpdatePost :one
UPDATE posts
SET
    title = $2,
    slug = $3,
    content = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, user_id, slug, title, content, created_at, updated_at;

-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1;

-- name: ListPosts :many
SELECT id, user_id, slug, title, content, created_at, updated_at
FROM posts
ORDER BY created_at DESC;


-- name: GetPostBySlug :one
SELECT id, user_id, slug, title, content, created_at, updated_at
FROM posts
WHERE slug = $1;

