BEGIN;

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    product_name VARCHAR(255) NOT NULL UNIQUE,
    product_address TEXT NOT NULL,
    product_time TIME NOT NULL,
    product_date DATE NOT NULL,
    product_price NUMERIC(13, 2) NOT NULL,
    product_description TEXT NOT NULL,
    product_quantity INT NOT NULL,
    product_type VARCHAR(20) CHECK (product_type IN ('available', 'unavailable')) NOT NULL,
    product_status VARCHAR(20) CHECK (product_status IN ('unpaid', 'pending', 'rejected', 'accepted')) NOT NULL,
    product_category VARCHAR(255) NOT NULL,
    user_id INT NOT NULL,
    CONSTRAINT fk_user_products FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION update_updated_at_column_products()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = CURRENT_TIMESTAMP;
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at_products
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column_products();

COMMIT;
