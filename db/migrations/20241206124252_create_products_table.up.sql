BEGIN;

CREATE TABLE IF NOT EXISTS products (
    id INT PRIMARY KEY AUTO_INCREMENT,
    product_name VARCHAR(255) NOT NULL,
    product_address TEXT NOT NULL,
    product_time TIME NOT NULL,
    product_date DATE NOT NULL,
    product_price DECIMAL(13, 2) NOT NULL,
    product_description TEXT NOT NULL,
    product_status ENUM('pending', 'ditolak', 'diterima') NOT NULL,
    product_category VARCHAR(255) NOT NULL,
    user_id INT NOT NULL,
    CONSTRAINT fk_user_products FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

COMMIT;