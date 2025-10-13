-- name: UpdatePost :one
UPDATE posts
SET title = $2, content = $3, status = $4, updated_at = $5
WHERE id = $1
RETURNING id, title, content, author_id, status, created_at, updated_at;