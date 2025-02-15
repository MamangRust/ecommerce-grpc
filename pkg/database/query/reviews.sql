-- name: GetReviews :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM reviews
WHERE deleted_at IS NULL
AND ($1::TEXT IS NULL OR review_id::TEXT ILIKE '%' || $1 || '%' OR name ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- name: GetReviewsActive :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM reviews
WHERE deleted_at IS NULL
AND ($1::TEXT IS NULL OR review_id::TEXT ILIKE '%' || $1 || '%' OR name ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- name: GetReviewsTrashed :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM reviews
WHERE deleted_at IS NOT NULL
AND ($1::TEXT IS NULL OR review_id::TEXT ILIKE '%' || $1 || '%' OR name ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- name: GetReviewsByProductID :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM reviews
WHERE deleted_at IS NULL
AND product_id = $1
AND ($2::TEXT IS NULL OR review_id::TEXT ILIKE '%' || $1 || '%' OR name ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;


-- name: GetReviewByID :one
SELECT *
FROM reviews
WHERE review_id = $1
  AND deleted_at IS NULL;

-- name: CreateReview :one
INSERT INTO reviews (
    user_id, product_id, name, comment, rating
) VALUES ($1, $2, $3, $4, $5)
RETURNING *;


-- name: UpdateReview :one
UPDATE reviews
SET 
    name = $2,
    comment = $3,
    rating = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE review_id = $1
AND deleted_at IS NULL
RETURNING *;


-- name: TrashReview :one
UPDATE reviews
SET deleted_at = CURRENT_TIMESTAMP
WHERE review_id = $1
AND deleted_at IS NULL
RETURNING *;


-- name: RestoreReview :one
UPDATE reviews
SET deleted_at = NULL
WHERE review_id = $1
AND deleted_at IS NOT NULL
RETURNING *;


-- name: DeleteReviewPermanently :exec
DELETE FROM reviews WHERE review_id = $1 AND deleted_at IS NOT NULL;


-- name: RestoreAllReviews :exec
UPDATE reviews
SET deleted_at = NULL
WHERE deleted_at IS NOT NULL;


-- name: DeleteAllPermanentReviews :exec
DELETE FROM reviews
WHERE deleted_at IS NOT NULL;
