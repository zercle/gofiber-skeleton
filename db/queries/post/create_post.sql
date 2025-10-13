-- name: CreatePost :one
INSERT INTO posts (id, title, content, author_id, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, title, content, author_id, status, created_at, updated_at;