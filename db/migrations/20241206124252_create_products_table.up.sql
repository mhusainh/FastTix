BEGIN;

CREATE TABLE IF NOT EXISTS products (
    id INT PRIMARY KEY AUTO_INCREMENT,
    product_name VARCHAR(255) NOT NULL UNIQUE,
    product_address TEXT NOT NULL,
    product_time TIME NOT NULL,
    product_date DATE NOT NULL,
    product_price DECIMAL(13, 2) NOT NULL,
    product_description TEXT NOT NULL,
    product_quantity INT NOT NULL,
    product_type ENUM('available', 'unavailable') NOT NULL,
    product_status ENUM('unpaid','pending', 'rejected', 'accepted') NOT NULL,
    product_category VARCHAR(255) NOT NULL,
    order_id VARCHAR(255) NOT NULL UNIQUE,
    user_id INT NOT NULL,
    CONSTRAINT fk_user_products FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

COMMIT;
