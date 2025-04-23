-- +goose Up
-- +goose StatementBegin
CREATE TABLE "orders" (
    "order_id" SERIAL PRIMARY KEY,
    "user_id" INT NOT NULL REFERENCES "users" ("user_id"),
    "merchant_id" INT NOT NULL REFERENCES "merchants" ("merchant_id"),
    "total_price" INT NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL
);
CREATE INDEX "idx_orders_user_id" ON "orders"("user_id");
CREATE INDEX "idx_orders_merchant_id" ON "orders"("merchant_id");
CREATE INDEX "idx_orders_phone" ON "orders"("phone");
CREATE INDEX "idx_orders_email" ON "orders"("email");
CREATE INDEX "idx_orders_courier" ON "orders"("courier");
CREATE INDEX "idx_orders_shipping_method" ON "orders"("shipping_method");
CREATE INDEX "idx_orders_total_price" ON "orders"("total_price");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "idx_orders_user_id";
DROP INDEX IF EXISTS "idx_orders_merchant_id";
DROP INDEX IF EXISTS "idx_orders_phone";
DROP INDEX IF EXISTS "idx_orders_email";
DROP INDEX IF EXISTS "idx_orders_courier";
DROP INDEX IF EXISTS "idx_orders_shipping_method";
DROP INDEX IF EXISTS "idx_orders_total_price";

DROP TABLE IF EXISTS "orders";
-- +goose StatementEnd
