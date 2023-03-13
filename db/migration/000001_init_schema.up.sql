-- Create price_data Table
CREATE TABLE "price_data" (
    "id" bigserial PRIMARY KEY,
    "symbol" varchar NOT NULL,
    "open" DECIMAL(20,10) NOT NULL,
    "close" DECIMAL(20,10) NOT NULL,
    "low" DECIMAL(20,10) NOT NULL,
    "high" DECIMAL(20,10) NOT NULL,
    "unix" bigint NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- Create Indexes
CREATE INDEX idx_price_data_symbol ON "price_data" ("symbol");
CREATE INDEX idx_price_data_unix ON "price_data" ("unix");
CREATE INDEX idx_price_data_created_at ON "price_data" ("created_at");
