-- name: GetUserByID :one
SELECT id, email, password_hash, first_name, last_name, is_active, created_at, updated_at, deleted_at
FROM users
WHERE id = ? AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, first_name, last_name, is_active, created_at, updated_at, deleted_at
FROM users
WHERE email = ? AND deleted_at IS NULL;

-- name: ListUsers :many
SELECT id, email, password_hash, first_name, last_name, is_active, created_at, updated_at, deleted_at
FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CreateUser :exec
INSERT INTO users (id, email, password_hash, first_name, last_name, is_active, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW());

-- name: UpdateUser :exec
UPDATE users
SET email = ?, password_hash = ?, first_name = ?, last_name = ?, is_active = ?, updated_at = NOW()
WHERE id = ? AND deleted_at IS NULL;

-- name: DeleteUser :exec
UPDATE users
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = ? AND deleted_at IS NULL;

-- name: HardDeleteUser :exec
DELETE FROM users
WHERE id = ?;

-- name: CountUsers :one
SELECT COUNT(*) as count
FROM users
WHERE deleted_at IS NULL;
