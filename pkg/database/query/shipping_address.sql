-- name: GetShippingAddress :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM shipping_addresses
WHERE deleted_at IS NULL
AND ($1::TEXT IS NULL OR shipping_address_id::TEXT ILIKE '%' || $1 || '%' OR alamat ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- name: GetShippingAddressActive :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM shipping_addresses
WHERE deleted_at IS NULL
AND ($1::TEXT IS NULL OR shipping_address_id::TEXT ILIKE '%' || $1 || '%' OR alamat ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- name: GetShippingAddressTrashed :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM shipping_addresses
WHERE deleted_at IS NOT NULL
AND ($1::TEXT IS NULL OR shipping_address_id::TEXT ILIKE '%' || $1 || '%' OR alamat ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- name: GetShippingByID :one
SELECT *
FROM shipping_addresses
WHERE shipping_address_id = $1
AND deleted_at IS NULL;

-- name: GetShippingAddressByOrderID :one
SELECT *
FROM shipping_addresses
WHERE order_id = $1
AND deleted_at IS NULL;


-- name: CreateShippingAddress :one
INSERT INTO shipping_addresses (
    order_id, alamat, provinsi, negara, kota, courier, shipping_method, shipping_cost
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;



-- name: UpdateShippingAddress :one
UPDATE shipping_addresses
SET 
    alamat = $2,
    provinsi = $3,
    negara = $4,
    kota = $5,
    courier = $6,
    shipping_method = $7,
    shipping_cost = $8,
    updated_at = CURRENT_TIMESTAMP
WHERE shipping_address_id = $1
AND deleted_at IS NULL
RETURNING *;



-- name: TrashShippingAddress :one
UPDATE shipping_addresses
SET deleted_at = CURRENT_TIMESTAMP
WHERE shipping_address_id = $1
AND deleted_at IS NULL
RETURNING *;


-- name: RestoreShippingAddress :one
UPDATE shipping_addresses
SET deleted_at = NULL
WHERE shipping_address_id = $1
AND deleted_at IS NOT NULL
RETURNING *;


-- name: DeleteShippingAddressPermanently :exec
DELETE FROM shipping_addresses WHERE shipping_address_id = $1 AND deleted_at IS NOT NULL;


-- name: RestoreAllShippingAddress :exec
UPDATE shipping_addresses
SET deleted_at = NULL
WHERE deleted_at IS NOT NULL;


-- name: DeleteAllPermanentShippingAddress :exec
DELETE FROM shipping_addresses
WHERE deleted_at IS NOT NULL;
