-- name: GetPostByID :one
SELECT id, title, content, author_id, status, created_at, updated_at
FROM posts
WHERE id = $1;