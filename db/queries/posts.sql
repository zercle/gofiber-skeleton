-- name: CreatePost :one
INSERT INTO posts (user_id, title, content, slug)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPostByID :one
SELECT 
    p.*,
    u.email as author_email,
    u.first_name as author_first_name,
    u.last_name as author_last_name
FROM posts p
JOIN users u ON p.user_id = u.id
WHERE p.id = $1;

-- name: GetPostBySlug :one
SELECT 
    p.*,
    u.email as author_email,
    u.first_name as author_first_name,
    u.last_name as author_last_name
FROM posts p
JOIN users u ON p.user_id = u.id
WHERE p.slug = $1;

-- name: ListPosts :many
SELECT 
    p.*,
    u.email as author_email,
    u.first_name as author_first_name,
    u.last_name as author_last_name
FROM posts p
JOIN users u ON p.user_id = u.id
WHERE ($3::boolean IS NULL OR p.is_published = $3)
ORDER BY p.created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListPostsByUser :many
SELECT * FROM posts
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdatePost :one
UPDATE posts
SET 
    title = COALESCE($2, title),
    content = COALESCE($3, content),
    slug = COALESCE($4, slug),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: PublishPost :one
UPDATE posts
SET 
    is_published = true,
    published_at = NOW(),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UnpublishPost :one
UPDATE posts
SET 
    is_published = false,
    published_at = NULL,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1;