-- name: CreateOrder :one
INSERT INTO orders (user_id, total_amount, shipping_address)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders
WHERE id = $1;

-- name: GetOrders :many
SELECT * FROM orders
ORDER BY created_at DESC;

-- name: GetUserOrders :many
SELECT * FROM orders
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateOrderStatus :one
UPDATE orders
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, unit_price, subtotal)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetOrderItems :many
SELECT * FROM order_items
WHERE order_id = $1;