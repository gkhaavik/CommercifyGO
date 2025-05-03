-- Drop indexes
DROP INDEX IF EXISTS idx_exchange_rate_history_base_currency;
DROP INDEX IF EXISTS idx_exchange_rate_history_date;
DROP INDEX IF EXISTS idx_currencies_is_enabled;

-- Drop the exchange rate history table
DROP TABLE IF EXISTS exchange_rate_history;

-- Drop the currencies table
DROP TABLE IF EXISTS currencies;

-- Remove currency columns from various tables
ALTER TABLE shipping_rates DROP COLUMN IF EXISTS currency;
ALTER TABLE discounts DROP COLUMN IF EXISTS currency;
ALTER TABLE orders DROP COLUMN IF EXISTS currency;
ALTER TABLE products DROP COLUMN IF EXISTS currency;