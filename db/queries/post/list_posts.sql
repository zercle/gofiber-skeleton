-- name: ListPosts :many
SELECT id, title, content, author_id, status, created_at, updated_at
FROM posts
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListPostsByAuthor :many
SELECT id, title, content, author_id, status, created_at, updated_at
FROM posts
WHERE author_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;