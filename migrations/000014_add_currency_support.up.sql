-- Add currency support to products table
ALTER TABLE products
ADD COLUMN currency VARCHAR(3) NOT NULL DEFAULT 'USD';

-- Add currency support to orders table
ALTER TABLE orders
ADD COLUMN currency VARCHAR(3) NOT NULL DEFAULT 'USD';

-- Add currency support to discounts table
ALTER TABLE discounts
ADD COLUMN currency VARCHAR(3) NOT NULL DEFAULT 'USD';

-- Add currency support to shipping rates table
ALTER TABLE shipping_rates
ADD COLUMN currency VARCHAR(3) NOT NULL DEFAULT 'USD';

-- Create a currencies table to store configured currencies
CREATE TABLE currencies (
    code VARCHAR(3) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    symbol VARCHAR(5) NOT NULL,
    precision INTEGER NOT NULL DEFAULT 2,
    exchange_rate DECIMAL(12, 6) NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    is_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Insert default currencies
INSERT INTO currencies (code, name, symbol, precision, exchange_rate, is_default, is_enabled)
VALUES
    ('USD', 'US Dollar', '$', 2, 1.0, TRUE, TRUE),
    ('EUR', 'Euro', '€', 2, 0.85, FALSE, TRUE),
    ('GBP', 'British Pound', '£', 2, 0.75, FALSE, TRUE),
    ('JPY', 'Japanese Yen', '¥', 0, 110.0, FALSE, TRUE),
    ('CAD', 'Canadian Dollar', '$', 2, 1.25, FALSE, TRUE);

-- Create an exchange_rate_history table to track rate changes over time
CREATE TABLE exchange_rate_history (
    id SERIAL PRIMARY KEY,
    base_currency VARCHAR(3) NOT NULL REFERENCES currencies(code),
    target_currency VARCHAR(3) NOT NULL REFERENCES currencies(code),
    rate DECIMAL(12, 6) NOT NULL,
    date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE (base_currency, target_currency, date)
);

-- Add indexes for performance
CREATE INDEX idx_currencies_is_enabled ON currencies(is_enabled);
CREATE INDEX idx_exchange_rate_history_date ON exchange_rate_history(date);
CREATE INDEX idx_exchange_rate_history_base_currency ON exchange_rate_history(base_currency);