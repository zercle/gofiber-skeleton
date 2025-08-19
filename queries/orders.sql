-- name: CreateOrder :one
INSERT INTO orders (
    user_id,
    total_price,
    status
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetOrderByID :one
SELECT * FROM orders
WHERE id = $1 LIMIT 1;

-- name: ListOrders :many
SELECT * FROM orders
ORDER BY created_at DESC;

-- name: UpdateOrderStatus :one
UPDATE orders
SET status = $2
WHERE id = $1
RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (
    order_id,
    product_id,
    quantity,
    price
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: ListOrderItemsByOrderID :many
SELECT * FROM order_items
WHERE order_id = $1
ORDER BY created_at DESC;