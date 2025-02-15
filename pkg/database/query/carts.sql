-- name: GetCarts :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM carts
WHERE deleted_at IS NULL
AND user_id = $1
AND ($2::TEXT IS NULL OR name ILIKE '%' || $2 || '%' OR price ILIKE '%' || $2 || '%')
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;


-- name: CreateCart :one
INSERT INTO "carts" ("user_id", "product_id", "name", "price", "image", "quantity", "weight")
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;


-- name: DeleteCart :exec
DELETE FROM "carts" WHERE "cart_id" = $1;


-- name: DeleteAllCart :exec
DELETE FROM "carts" WHERE "cart_id" = ANY($1::int[]);
