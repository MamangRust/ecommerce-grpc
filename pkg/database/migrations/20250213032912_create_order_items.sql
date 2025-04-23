-- +goose Up
-- +goose StatementBegin
CREATE TABLE "order_items" (
    "order_item_id" SERIAL PRIMARY KEY,
    "order_id" INT NOT NULL REFERENCES "orders" ("order_id"),
    "product_id" INT NOT NULL REFERENCES "products" ("product_id"),
    "name" VARCHAR(255) NOT NULL,
    "quantity" INT NOT NULL,
    "price" INT NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL
);
CREATE INDEX "idx_order_items_order_id" ON "order_items"("order_id");
CREATE INDEX "idx_order_items_name" ON "order_items"("name");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "idx_order_items_order_id";
DROP INDEX IF EXISTS "idx_order_items_name";

DROP TABLE IF EXISTS "order_items"
-- +goose StatementEnd
