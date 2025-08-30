-- name: CreateProduct :one
INSERT INTO products (name, description, price, stock, image_url)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetProductByID :one
SELECT * FROM products WHERE id = $1;

-- name: GetAllProducts :many
SELECT * FROM products ORDER BY created_at DESC;

-- name: UpdateProduct :one
UPDATE products 
SET name = $2, description = $3, price = $4, stock = $5, image_url = $6, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;

-- name: UpdateProductStock :one
UPDATE products 
SET stock = stock + $2, updated_at = NOW()
WHERE id = $1
RETURNING *;