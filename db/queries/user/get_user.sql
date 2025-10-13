-- name: GetUserByID :one
SELECT id, email, password_hash, full_name, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, full_name, created_at, updated_at
FROM users
WHERE email = $1;