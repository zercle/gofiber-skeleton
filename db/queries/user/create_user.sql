-- name: CreateUser :one
INSERT INTO users (id, email, password_hash, full_name, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, email, password_hash, full_name, created_at, updated_at;