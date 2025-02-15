-- name: GetOrders :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM orders
WHERE deleted_at IS NULL
AND ($1::TEXT IS NULL OR order_id::TEXT ILIKE '%' || $1 || '%' OR total_price::TEXT ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- name: GetOrdersActive :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM orders
WHERE deleted_at IS NULL
AND ($1::TEXT IS NULL OR order_id::TEXT ILIKE '%' || $1 || '%' OR total_price::TEXT ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- name: GetOrdersTrashed :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM orders
WHERE deleted_at IS NOT NULL
AND ($1::TEXT IS NULL OR order_id::TEXT ILIKE '%' || $1 || '%' OR total_price::TEXT ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- name: GetOrdersByMerchant :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM orders
WHERE 
    deleted_at IS NULL
    AND ($1::TEXT IS NULL OR order_id::TEXT ILIKE '%' || $1 || '%' OR total_price::TEXT ILIKE '%' || $1 || '%')
    AND ($4::INT IS NULL OR merchant_id = $4)
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;



-- name: CreateOrder :one
INSERT INTO orders (merchant_id, user_id, total_price)
VALUES ($1, $2, $3)
RETURNING *;


-- name: GetOrderByID :one
SELECT *
FROM orders
WHERE order_id = $1
AND deleted_at IS NULL;


-- name: UpdateOrder :one
UPDATE orders
SET total_price = $2,
    user_id = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE order_id = $1
  AND deleted_at IS NULL
  RETURNING *;  


-- name: TrashOrder :one
UPDATE orders
SET deleted_at = CURRENT_TIMESTAMP
WHERE order_id = $1
AND deleted_at IS NULL
RETURNING *;


-- name: RestoreOrder :one
UPDATE orders
SET deleted_at = NULL
WHERE order_id = $1
AND deleted_at IS NOT NULL
RETURNING *;


-- name: DeleteOrderPermanently :exec
DELETE FROM orders WHERE order_id = $1 AND deleted_at IS NOT NULL;


-- name: RestoreAllOrders :exec
UPDATE orders
SET deleted_at = NULL
WHERE deleted_at IS NOT NULL;


-- name: DeleteAllPermanentOrders :exec
DELETE FROM orders
WHERE deleted_at IS NOT NULL;
