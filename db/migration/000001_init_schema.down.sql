-- Drop indexes
DROP INDEX IF EXISTS idx_price_data_symbol;
DROP INDEX IF EXISTS idx_price_data_unix;
DROP INDEX IF EXISTS idx_price_data_created_at;

-- Drop Tables
DROP TABLE "price_data";