-- name: CreatePost :one
INSERT INTO posts (id, title, content, status, user_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetPostByID :one
SELECT * FROM posts WHERE id = $1;

-- name: GetPostsByUserID :many
SELECT * FROM posts WHERE user_id = $1 ORDER BY created_at DESC;

-- name: GetAllPosts :many
SELECT p.*, u.full_name, u.email
FROM posts p
JOIN users u ON p.user_id = u.id
ORDER BY p.created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdatePost :one
UPDATE posts
SET title = $2, content = $3, status = $4, updated_at = $5
WHERE id = $1
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts WHERE id = $1;

-- name: GetPostsWithAuthor :many
SELECT p.*, u.full_name, u.email
FROM posts p
JOIN users u ON p.user_id = u.id
WHERE p.status = 'published'
ORDER BY p.created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetUserPostStats :one
SELECT
    COUNT(*) as total_posts,
    COUNT(CASE WHEN status = 'published' THEN 1 END) as published_posts,
    COUNT(CASE WHEN status = 'draft' THEN 1 END) as draft_posts,
    MAX(created_at) as last_post_date
FROM posts
WHERE user_id = $1;

-- name: GetPublishedPostsByUserID :many
SELECT * FROM posts
WHERE user_id = $1 AND status = 'published'
ORDER BY created_at DESC;

-- name: GetDraftPostsByUserID :many
SELECT * FROM posts
WHERE user_id = $1 AND status = 'draft'
ORDER BY created_at DESC;

-- name: CountPostsByStatus :one
SELECT
    status,
    COUNT(*) as count
FROM posts
GROUP BY status;

-- name: SearchPosts :many
SELECT p.*, u.full_name, u.email
FROM posts p
JOIN users u ON p.user_id = u.id
WHERE (p.title ILIKE $1 OR p.content ILIKE $1)
  AND p.status = $2
ORDER BY p.created_at DESC
LIMIT $3 OFFSET $4;