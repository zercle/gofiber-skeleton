-- name: CreateUser :exec
INSERT INTO users (
    username,
    email,
    password_hash,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC;