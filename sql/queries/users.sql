-- name: CreateUser :one
INSERT INTO users(id, 
    name, 
    username,
    email, 
    password, 
    role, 
    suspended_until) VALUES
    (
        $1, 
        $2, 
        $3, 
        $4, 
        $5, 
        $6,
        $7
    )
    RETURNING id, name, email, role, suspended_until, created_at, updated_at,username;

-- name: GetUserByEmail :one
SELECT id, 
    name, 
    username,
    email, 
    password, 
    role, 
    suspended_until, 
    created_at, 
    updated_at FROM users WHERE email = $1;

-- name: GetUserById :one
SELECT id, 
    name, 
    username,
    email, 
    password, 
    role, 
    suspended_until, 
    created_at, 
    updated_at FROM users WHERE id = $1;


-- name: GetUserByUsername :one
SELECT id, 
    name, 
    username,
    email, 
    password, 
    role, 
    suspended_until, 
    created_at, 
    updated_at FROM users WHERE username = $1;