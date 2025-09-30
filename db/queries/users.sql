-- name: CreateUser :one
INSERT INTO users (
    username,
    email,
    password_hash
) VALUES (
    $1, $2, $3
)
RETURNING id, username, email, password_hash, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, username, email, password_hash, created_at, updated_at
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT id, username, email, password_hash, created_at, updated_at
FROM users
WHERE email = $1
LIMIT 1;

-- name: GetUserByUsername :one
SELECT id, username, email, password_hash, created_at, updated_at
FROM users
WHERE username = $1
LIMIT 1;

-- name: ListUsers :many
SELECT id, username, email, password_hash, created_at, updated_at
FROM users
ORDER BY created_at DESC;