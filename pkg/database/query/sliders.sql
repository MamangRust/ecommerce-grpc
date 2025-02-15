-- name: GetSliders :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM sliders
WHERE deleted_at IS NULL
AND ($1::TEXT IS NULL OR name ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- name: GetSlidersActive :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM sliders
WHERE deleted_at IS NULL
AND ($1::TEXT IS NULL OR name ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- name: GetSlidersTrashed :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM sliders
WHERE deleted_at IS NOT NULL
AND ($1::TEXT IS NULL OR name ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- name: CreateSlider :one
INSERT INTO sliders (name, image)
VALUES ($1, $2)
RETURNING *;


-- name: GetSliderByID :one
SELECT *
FROM sliders
WHERE slider_id = $1
AND deleted_at IS NULL;


-- name: UpdateSlider :one
UPDATE sliders
SET name = $2,
    image = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE slider_id = $1
AND deleted_at IS NULL
RETURNING *;


-- name: TrashSlider :one
UPDATE sliders
SET deleted_at = CURRENT_TIMESTAMP
WHERE slider_id = $1
AND deleted_at IS NULL
RETURNING *;



-- name: RestoreSlider :one
UPDATE sliders
SET deleted_at = NULL
WHERE slider_id = $1
AND deleted_at IS NOT NULL
RETURNING *;


-- name: DeleteSliderPermanently :exec
DELETE FROM sliders WHERE slider_id = $1 AND deleted_at IS NOT NULL;


-- name: RestoreAllSliders :exec
UPDATE sliders
SET deleted_at = NULL
WHERE deleted_at IS NOT NULL;


-- name: DeleteAllPermanentSliders :exec
DELETE FROM sliders
WHERE deleted_at IS NOT NULL;
