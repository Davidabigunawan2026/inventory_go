
-- =========================================
-- INVENTORY GO
-- Migration 001: Create Products Table
-- =========================================

CREATE TABLE IF NOT EXISTS products (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    sku VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(150) NOT NULL,
    description TEXT,

    category VARCHAR(100),
    unit VARCHAR(30) NOT NULL DEFAULT 'PCS',

    min_stock NUMERIC(14,2) NOT NULL DEFAULT 0,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_products_min_stock
        CHECK (min_stock >= 0)
);

-- Index untuk pencarian nama dan kategori
CREATE INDEX IF NOT EXISTS idx_products_name
    ON products(name);

CREATE INDEX IF NOT EXISTS idx_products_category
    ON products(category);
