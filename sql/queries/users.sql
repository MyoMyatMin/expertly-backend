-- name: CreateUser :one
INSERT INTO users(id, 
    name, 
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
        $6
    )
    RETURNING id, name, email, role, suspended_until, created_at, updated_at;

-- name: GetUserByEmail :one
SELECT id, 
    name, 
    email, 
    password, 
    role, 
    suspended_until, 
    created_at, 
    updated_at FROM users WHERE email = $1;

-- name: GetUserById :one
SELECT id, 
    name, 
    email, 
    password, 
    role, 
    suspended_until, 
    created_at, 
    updated_at FROM users WHERE id = $1;