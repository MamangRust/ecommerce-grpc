-- Get Categories with Pagination and Total Count
-- name: GetCategories :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM categories
WHERE deleted_at IS NULL
AND ($1::TEXT IS NULL OR name ILIKE '%' || $1 || '%' OR slug_category ILIKE '%' || $1 || '%') --
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- Get Active Categories with Pagination and Total Count
-- name: GetCategoriesActive :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM categories
WHERE deleted_at IS NULL
AND ($1::TEXT IS NULL OR name ILIKE '%' || $1 || '%' OR slug_category ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- Get Trashed Categories with Pagination and Total Count
-- name: GetCategoriesTrashed :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM categories
WHERE deleted_at IS NOT NULL
AND ($1::TEXT IS NULL OR name ILIKE '%' || $1 || '%' OR slug_category ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;



-- name: CreateCategory :one
INSERT INTO categories (name, description, slug_category, image_category)
VALUES ($1, $2, $3, $4)
  RETURNING *;

-- name: GetCategoryByID :one
SELECT *
FROM categories
WHERE category_id = $1
  AND deleted_at IS NULL;


-- name: UpdateCategory :one
UPDATE categories
SET name = $2,
    description = $3,
    slug_category = $4,
    image_category = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE category_id = $1
  AND deleted_at IS NULL
  RETURNING *;

-- Trash Category
-- name: TrashCategory :one
UPDATE categories
SET
    deleted_at = current_timestamp
WHERE
    category_id = $1
    AND deleted_at IS NULL
    RETURNING *;    


-- Restore Trashed Category
-- name: RestoreCategory :one
UPDATE categories
SET
    deleted_at = NULL
WHERE
    category_id = $1
    AND deleted_at IS NOT NULL
  RETURNING *;


-- Delete Category Permanently
-- name: DeleteCategoryPermanently :exec
DELETE FROM categories WHERE category_id = $1 AND deleted_at IS NOT NULL;


-- Restore All Trashed Category
-- name: RestoreAllCategories :exec
UPDATE categories
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;


-- Delete All Trashed Category Permanently
-- name: DeleteAllPermanentCategories :exec
DELETE FROM categories
WHERE
    deleted_at IS NOT NULL;