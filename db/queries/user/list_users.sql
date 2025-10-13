-- name: ListUsers :many
SELECT id, email, password_hash, full_name, created_at, updated_at
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;