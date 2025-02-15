-- name: GetProducts :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM products
WHERE deleted_at IS NULL
AND ($1::TEXT IS NULL 
       OR p.name ILIKE '%' || $1 || '%'
       OR p.description ILIKE '%' || $1 || '%'
       OR p.brand ILIKE '%' || $1 || '%'
       OR p.slug_product ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- Get Active Products with Pagination and Total Count
-- name: GetProductsActive :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM products
WHERE deleted_at IS NULL
AND ($1::TEXT IS NULL 
       OR p.name ILIKE '%' || $1 || '%'
       OR p.description ILIKE '%' || $1 || '%'
       OR p.brand ILIKE '%' || $1 || '%'
       OR p.slug_product ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- Get Trashed Products with Pagination and Total Count
-- name: GetProductsTrashed :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM products
WHERE deleted_at IS NOT NULL
AND ($1::TEXT IS NULL 
       OR p.name ILIKE '%' || $1 || '%'
       OR p.description ILIKE '%' || $1 || '%'
       OR p.brand ILIKE '%' || $1 || '%'
       OR p.slug_product ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;



-- name: GetProductsByMerchant :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM products
WHERE merchant_id = $1
AND deleted_at IS NULL
AND ($2::TEXT IS NULL 
       OR p.name ILIKE '%' || $2 || '%'
       OR p.description ILIKE '%' || $2 || '%'
       OR p.brand ILIKE '%' || $2 || '%'
       OR p.slug_product ILIKE '%' || $2 || '%'
       OR p.barcode ILIKE '%' || $2 || '%')
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;


-- Get Products by Category Name with Filters
-- name: GetProductsByCategoryName :many
SELECT
    p.*,
    COUNT(*) OVER() AS total_count
FROM products p
JOIN categories c ON p.category_id = c.category_id
WHERE c.deleted_at IS NULL
  AND c.name = $1 
  AND ($2::TEXT IS NULL 
       OR p.name ILIKE '%' || $2 || '%'
       OR p.description ILIKE '%' || $2 || '%'
       OR p.brand ILIKE '%' || $2 || '%'
       OR p.slug_product ILIKE '%' || $2 || '%')
ORDER BY p.created_at DESC
LIMIT $3 OFFSET $4;




-- name: CreateProduct :one
INSERT INTO products (merchant_id, category_id, name, description, price, count_in_stock, brand, weight, rating, slug_product, image_product)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: GetProductByID :one
SELECT *
FROM products
WHERE product_id = $1
  AND deleted_at IS NULL;


-- name: UpdateProduct :one
UPDATE products
SET category_id = $2,
    name = $3,
    description = $4,
    price = $5,
    count_in_stock = $6,
    brand = $7,
    weight = $8,
    rating = $9,
    slug_product = $10,
    image_product = $11,
    updated_at = CURRENT_TIMESTAMP
WHERE product_id = $1
  AND deleted_at IS NULL
  RETURNING *;


-- name: UpdateProductCountStock :one
UPDATE products
SET count_in_stock = $2
WHERE product_id = $1
    AND deleted_at IS NULL
RETURNING *;

-- Trash Product
-- name: TrashProduct :one
UPDATE products
SET
    deleted_at = current_timestamp
WHERE
    product_id = $1
    AND deleted_at IS NULL
    RETURNING *;    


-- Restore Trashed Product
-- name: RestoreProduct :one
UPDATE products
SET
    deleted_at = NULL
WHERE
    product_id = $1
    AND deleted_at IS NOT NULL
  RETURNING *;


-- Delete Product Permanently
-- name: DeleteProductPermanently :exec
DELETE FROM products WHERE product_id = $1 AND deleted_at IS NOT NULL;


-- Restore All Trashed Product
-- name: RestoreAllProducts :exec
UPDATE products
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;


-- Delete All Trashed Product Permanently
-- name: DeleteAllPermanentProducts :exec
DELETE FROM products
WHERE
    deleted_at IS NOT NULL;
