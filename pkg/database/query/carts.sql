-- GetCarts: Retrieves paginated list of user's active cart items with search capability
-- Purpose: Display cart items for e-commerce checkout interface
-- Parameters:
--   $1: user_id - ID of the user whose cart to retrieve
--   $2: search_term - Optional text to filter items by name or price (NULL for no filter)
--   $3: limit - Maximum number of records to return
--   $4: offset - Number of records to skip for pagination
-- Returns:
--   All cart fields plus total_count of matching records
-- Business Logic:
--   - Excludes soft-deleted items (deleted_at IS NULL)
--   - Filters by specific user only
--   - Supports partial text matching on name and price fields (case-insensitive)
--   - Returns newest items first (created_at DESC)
--   - Provides total_count for pagination calculations
-- name: GetCarts :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM carts
WHERE deleted_at IS NULL
  AND user_id = $1
  AND (
    $2::TEXT IS NULL OR 
    name ILIKE '%' || $2 || '%' OR 
    price::TEXT ILIKE '%' || $2 || '%'
  )
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;



-- CreateCart: Adds a new product to user's shopping cart
-- Purpose: Add items to cart during product browsing
-- Parameters:
--   $1: user_id - ID of the user adding the item
--   $2: product_id - ID of the product being added
--   $3: name - Display name of the product
--   $4: price - Current price of the product
--   $5: image - URL of product image
--   $6: quantity - Number of units being added
--   $7: weight - Total weight of the item(s) for shipping
-- Returns:
--   The complete created cart record
-- Business Logic:
--   - Creates a new cart entry with all product details
--   - Requires all product information for historical accuracy
--   - Returns the full record for immediate UI update
-- name: CreateCart :one
INSERT INTO "carts" ("user_id", "product_id", "name", "price", "image", "quantity", "weight")
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;


-- DeleteCart: Removes a single item from user's cart
-- Purpose: Remove unwanted items during checkout process
-- Parameters:
--   $1: cart_id - ID of the specific cart item to remove
-- Returns:
--   Nothing (exec-only)
-- Business Logic:
--   - Permanently deletes the specified cart item
--   - Used for individual item removal
--   - No return value needed (UI refreshes cart after operation)
-- name: DeleteCart :exec
DELETE FROM "carts" WHERE "cart_id" = $1;


-- DeleteAllCart: Batch removes multiple cart items
-- Purpose: Clear selected items or reset cart
-- Parameters:
--   $1: cart_ids - Array of cart item IDs to remove
-- Returns:
--   Nothing (exec-only)
-- Business Logic:
--   - Efficiently deletes multiple items in single operation
--   - Uses array parameter for batch processing
--   - Typically used for "Clear Cart" functionality
-- name: DeleteAllCart :exec
DELETE FROM "carts" WHERE "cart_id" = ANY($1::int[]);