-- name: GetMerchants :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM merchants
WHERE deleted_at IS NULL
AND ($1::TEXT IS NULL OR name ILIKE '%' || $1 || '%' OR contact_email ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- Get Active Merchants with Pagination and Total Count
-- name: GetMerchantsActive :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM merchants
WHERE deleted_at IS NULL
AND ($1::TEXT IS NULL OR name ILIKE '%' || $1 || '%' OR contact_email ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- Get Trashed Merchants with Pagination and Total Count
-- name: GetMerchantsTrashed :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM merchants
WHERE deleted_at IS NOT NULL
AND ($1::TEXT IS NULL OR name ILIKE '%' || $1 || '%' OR contact_email ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- Create Merchant
-- name: CreateMerchant :one
INSERT INTO merchants (user_id, name, description, address, contact_email, contact_phone, status)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetMerchantByID :one
SELECT *
FROM merchants
WHERE merchant_id = $1
  AND deleted_at IS NULL;

-- Update Merchant
-- name: UpdateMerchant :one
UPDATE merchants
SET name = $2,
    description = $3,
    address = $4,
    contact_email = $5,
    contact_phone = $6,
    status = $7,
    updated_at = CURRENT_TIMESTAMP
WHERE merchant_id = $1
  AND deleted_at IS NULL
  RETURNING *;


-- Trash Merchant
-- name: TrashMerchant :one
UPDATE merchants
SET
    deleted_at = current_timestamp
WHERE
    merchant_id = $1
    AND deleted_at IS NULL
    RETURNING *;    


-- Restore Trashed Merchant
-- name: RestoreMerchant :one
UPDATE merchants
SET
    deleted_at = NULL
WHERE
    merchant_id = $1
    AND deleted_at IS NOT NULL
  RETURNING *;


-- Delete Merchant Permanently
-- name: DeleteMerchantPermanently :exec
DELETE FROM merchants WHERE merchant_id = $1 AND deleted_at IS NOT NULL;


-- Restore All Trashed Merchant
-- name: RestoreAllMerchants :exec
UPDATE merchants
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;


-- Delete All Trashed Merchant Permanently
-- name: DeleteAllPermanentMerchants :exec
DELETE FROM merchants
WHERE
    deleted_at IS NOT NULL;

