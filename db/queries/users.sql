-- name: CreateUser :one
INSERT INTO users (
    email,
    password_hash,
    full_name
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
    full_name = COALESCE($2, full_name),
    is_active = COALESCE($3, is_active),
    is_verified = COALESCE($4, is_verified)
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $2
WHERE id = $1;

-- name: VerifyUser :exec
UPDATE users
SET is_verified = true
WHERE id = $1;
