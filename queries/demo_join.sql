-- name: GetOrdersWithItemsAndProducts :many
SELECT
    o.id AS order_id,
    o.user_id,
    o.status AS order_status,
    o.total AS order_total,
    o.created_at AS order_created_at,
    oi.id AS order_item_id,
    oi.product_id,
    oi.quantity,
    oi.price AS item_price,
    p.name AS product_name,
    p.description AS product_description,
    p.price AS product_unit_price,
    p.stock AS product_stock,
    p.image_url AS product_image_url
FROM orders o
JOIN order_items oi ON o.id = oi.order_id
JOIN products p ON oi.product_id = p.id
ORDER BY o.created_at DESC;