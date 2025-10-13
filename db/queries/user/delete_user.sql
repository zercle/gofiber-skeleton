-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;