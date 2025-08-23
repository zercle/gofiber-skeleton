-- name: CreateOrder :one
INSERT INTO orders (user_id, status, total)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetOrderByID :one
SELECT * FROM orders WHERE id = $1;

-- name: GetOrdersByUserID :many
SELECT * FROM orders WHERE user_id = $1 ORDER BY created_at DESC;

-- name: GetAllOrders :many
SELECT * FROM orders ORDER BY created_at DESC;

-- name: UpdateOrderStatus :one
UPDATE orders 
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateOrder :one
UPDATE orders 
SET user_id = $2, status = $3, total = $4, updated_at = NOW()
WHERE id = $1
RETURNING *;