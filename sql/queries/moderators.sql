-- name: CreateModerator :one
INSERT INTO moderators(
    moderator_id, 
    name, 
    email, 
    password, 
    role, 
    created_by
) VALUES (
    $1, 
    $2, 
    $3, 
    $4, 
    $5, 
    $6
)
RETURNING moderator_id, name, email, role, created_at, updated_at;

-- name: GetModeratorByEmail :one
SELECT 
    moderator_id, 
    name, 
    email, 
    password,
    role, 
    created_at, 
    updated_at
FROM moderators
WHERE email = $1;

-- name: GetModeratorById :one
SELECT 
    moderator_id, 
    name, 
    email, 
    role, 
    created_at, 
    updated_at
FROM moderators
WHERE moderator_id = $1;

-- name: GetALLModerators :many
SELECT 
    moderator_id, 
    name, 
    email, 
    role, 
    created_at, 
    updated_at
FROM moderators
ORDER BY created_at DESC;



