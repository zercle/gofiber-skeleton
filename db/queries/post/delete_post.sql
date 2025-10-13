-- name: DeletePost :exec
DELETE FROM posts WHERE id = $1;