
-- name: CreateUser :one
INSERT INTO users(
    user_id, 
    name, 
    username,
    email, 
    password, 
    suspended_until
) VALUES (
    $1, 
    $2, 
    $3, 
    $4, 
    $5, 
    $6
)
RETURNING user_id, name, email, username, suspended_until, created_at, updated_at;

-- name: GetUserByEmail :one
SELECT 
    user_id, 
    name, 
    username,
    email, 
    password, 
    suspended_until, 
    created_at, 
    updated_at 
FROM users 
WHERE email = $1;

-- name: GetUserById :one
SELECT 
    user_id, 
    name, 
    username,
    email, 
    password, 
    suspended_until, 
    created_at, 
    updated_at 
FROM users 
WHERE user_id = $1;

-- name: GetUserByUsername :one
SELECT 
    user_id, 
    name, 
    username,
    email, 
    password, 
    suspended_until, 
    created_at, 
    updated_at 
FROM users 
WHERE username = $1;

-- name: GetIDbyUsername :one
SELECT user_id FROM users WHERE username = $1;


-- name: UpdateUser :exec
UPDATE users
SET 
    name = $2, 
    username = $3
WHERE user_id = $1
RETURNING user_id, name, email, username, suspended_until, created_at, updated_at;

