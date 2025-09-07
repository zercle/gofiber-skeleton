-- name: CreateUser :one
INSERT INTO users (email, password_hash, first_name, last_name)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 AND is_active = true;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 AND is_active = true;

-- name: UpdateUser :one
UPDATE users
SET 
    email = COALESCE($2, email),
    first_name = COALESCE($3, first_name),
    last_name = COALESCE($4, last_name),
    is_email_verified = COALESCE($5, is_email_verified),
    updated_at = NOW()
WHERE id = $1 AND is_active = true
RETURNING *;

-- name: UpdateUserLastLogin :exec
UPDATE users
SET last_login_at = NOW(), updated_at = NOW()
WHERE id = $1 AND is_active = true;

-- name: DeactivateUser :exec
UPDATE users
SET is_active = false, updated_at = NOW()
WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users
WHERE is_active = true
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;