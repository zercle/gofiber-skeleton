-- name: CreateUser :one
INSERT INTO users (
    username,
    email,
    password_hash,
    full_name
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
WHERE is_active = TRUE
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) FROM users
WHERE is_active = TRUE;

-- name: UpdateUser :one
UPDATE users
SET
    full_name = COALESCE(sqlc.narg('full_name'), full_name),
    email = COALESCE(sqlc.narg('email'), email),
    username = COALESCE(sqlc.narg('username'), username),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users
SET
    password_hash = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserLastLogin :exec
UPDATE users
SET
    last_login_at = NOW(),
    updated_at = NOW()
WHERE id = $1;

-- name: VerifyUser :exec
UPDATE users
SET
    is_verified = TRUE,
    updated_at = NOW()
WHERE id = $1;

-- name: DeactivateUser :exec
UPDATE users
SET
    is_active = FALSE,
    updated_at = NOW()
WHERE id = $1;

-- name: ActivateUser :exec
UPDATE users
SET
    is_active = TRUE,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
