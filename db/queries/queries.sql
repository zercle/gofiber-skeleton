-- name: CreateUser :one
INSERT INTO users (id, username, password) VALUES ($1, $2, $3) RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: CreateURL :one
INSERT INTO urls (id, user_id, short_code, long_url) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetURLByShortCode :one
SELECT * FROM urls WHERE short_code = $1;

-- name: GetURLsByUserID :many
SELECT * FROM urls WHERE user_id = $1;

-- name: DeleteURL :exec
DELETE FROM urls WHERE id = $1 AND user_id = $2;
