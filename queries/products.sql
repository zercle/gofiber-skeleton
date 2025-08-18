-- name: CreateProduct :one
INSERT INTO products (
    id, name, description, price, stock, image_url
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING id, name, description, price, stock, image_url, created_at, updated_at;

-- name: GetProductByID :one
SELECT id, name, description, price, stock, image_url, created_at, updated_at FROM products
WHERE id = $1 LIMIT 1;

-- name: GetAllProducts :many
SELECT id, name, description, price, stock, image_url, created_at, updated_at FROM products
ORDER BY name;

-- name: UpdateProduct :one
UPDATE products
SET
    name = $2,
    description = $3,
    price = $4,
    stock = $5,
    image_url = $6,
    updated_at = NOW()
WHERE id = $1
RETURNING id, name, description, price, stock, image_url, created_at, updated_at;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;