-- name: CreatePost :exec
INSERT INTO posts (
    user_id,
    thread_id,
    content,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: GetPostByID :one
SELECT * FROM posts WHERE id = $1 LIMIT 1;

-- name: ListPostsByUser :many
SELECT * FROM posts WHERE user_id = $1 ORDER BY created_at DESC;

-- name: UpdatePost :exec
UPDATE posts SET content = $2, updated_at = $3 WHERE id = $1;

-- name: DeletePost :exec
DELETE FROM posts WHERE id = $1;