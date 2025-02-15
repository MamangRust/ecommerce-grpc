-- Get Order Items  with Pagination and Total Count
-- name: GetOrderItems :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM order_items
WHERE deleted_at IS NULL
  AND ($1::TEXT IS NULL OR order_id::TEXT ILIKE '%' || $1 || '%' OR product_id::TEXT ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- Get Active Order Item with Pagination and Total Count
-- name: GetOrderItemsActive :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM order_items
WHERE deleted_at IS NULL
  AND ($1::TEXT IS NULL OR order_id::TEXT ILIKE '%' || $1 || '%' OR product_id::TEXT ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;



-- Get Trashed Orders Items with Pagination and Total Count
-- name: GetOrderItemsTrashed :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM order_items
WHERE deleted_at IS NOT NULL
  AND ($1::TEXT IS NULL OR order_id::TEXT ILIKE '%' || $1 || '%' OR product_id::TEXT ILIKE '%' || $1 || '%')
ORDER BY deleted_at DESC
LIMIT $2 OFFSET $3;



-- name: CalculateTotalPrice :one
SELECT COALESCE(SUM(quantity * price), 0)::int AS total_price
FROM order_items
WHERE order_id = $1 AND deleted_at IS NULL;



-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, price)
VALUES ($1, $2, $3, $4)
RETURNING *;


-- name: GetOrderItemsByOrder :many
SELECT *
FROM order_items
WHERE order_id = $1
  AND deleted_at IS NULL;

-- name: UpdateOrderItem :one
UPDATE order_items
SET quantity = $2,
    price = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE order_item_id = $1
  AND deleted_at IS NULL
  RETURNING *;

-- Trash Order Item
-- name: TrashOrderItem :one
UPDATE order_items
SET
    deleted_at = current_timestamp
WHERE
    order_id = $1
    AND deleted_at IS NULL
    RETURNING *;    


-- Restore Trashed Order Item
-- name: RestoreOrderItem :one
UPDATE order_items
SET
    deleted_at = NULL
WHERE
    order_id = $1
    AND deleted_at IS NOT NULL
  RETURNING *;


-- Delete Order Item Permanently
-- name: DeleteOrderItemPermanently :exec
DELETE FROM order_items WHERE order_id = $1 AND deleted_at IS NOT NULL;


-- Restore All Trashed Order Item
-- name: RestoreAllOrdersItem :exec
UPDATE order_items
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;


-- Delete All Trashed Order Item Permanently
-- name: DeleteAllPermanentOrdersItem :exec
DELETE FROM order_items 
WHERE
    deleted_at IS NOT NULL;
