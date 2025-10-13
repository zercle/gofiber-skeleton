-- name: UpdateUser :one
UPDATE users
SET email = $2, password_hash = $3, full_name = $4, updated_at = $5
WHERE id = $1
RETURNING id, email, password_hash, full_name, created_at, updated_at;