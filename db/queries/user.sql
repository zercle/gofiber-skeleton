-- name: CreateUser :one
INSERT INTO users (email, password_hash, full_name)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 AND is_active = true
LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 AND is_active = true
LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET email = $2, full_name = $3, updated_at = NOW()
WHERE id = $1 AND is_active = true
RETURNING *;

-- name: UpdateUserPassword :one
UPDATE users
SET password_hash = $2, updated_at = NOW()
WHERE id = $1 AND is_active = true
RETURNING *;

-- name: DeactivateUser :one
UPDATE users
SET is_active = false, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ListUsers :many
SELECT id, email, full_name, is_active, created_at, updated_at
FROM users
WHERE is_active = true
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;

-- name: UserExists :one
SELECT COUNT(*) FROM users
WHERE email = $1 AND is_active = true;