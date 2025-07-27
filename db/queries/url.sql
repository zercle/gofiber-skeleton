-- name: CreateURL :one
INSERT INTO urls (original_url, short_code, user_id, expires_at) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetURLByShortCode :one
SELECT * FROM urls WHERE short_code = $1;

-- name: GetURLsByUserID :many
SELECT * FROM urls WHERE user_id = $1;

-- name: UpdateURL :one
UPDATE urls SET original_url = $2 WHERE id = $1 RETURNING *;

-- name: DeleteURL :exec
DELETE FROM urls WHERE id = $1;

-- name: GetURLByID :one
SELECT * FROM urls WHERE id = $1;
