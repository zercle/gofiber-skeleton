-- name: CreatePost :one
INSERT INTO posts (
    user_id,
    thread_id,
    content
) VALUES (
    $1, $2, $3
)
RETURNING id, user_id, thread_id, content, created_at, updated_at;

-- name: GetPostByID :one
SELECT id, user_id, thread_id, content, created_at, updated_at
FROM posts
WHERE id = $1
LIMIT 1;

-- name: ListPostsByUser :many
SELECT id, user_id, thread_id, content, created_at, updated_at
FROM posts
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: ListPostsByThread :many
SELECT id, user_id, thread_id, content, created_at, updated_at
FROM posts
WHERE thread_id = $1
ORDER BY created_at ASC;

-- name: UpdatePost :one
UPDATE posts
SET content = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, thread_id, content, created_at, updated_at;

-- name: DeletePost :exec
DELETE FROM posts WHERE id = $1;