BEGIN;

CREATE TABLE IF NOT EXISTS transactions (
    id INT PRIMARY KEY AUTO_INCREMENT,
    transaction_status ENUM('pending', 'failed', 'success') NOT NULL,
    transaction_quantity INT NOT NULL,
    transaction_amount DECIMAL(13, 2) NOT NULL,
    tranaction_type ENUM('ticket', 'submission') NOT NULL,
    order_id VARCHAR(255) NOT NULL UNIQUE,
    user_id INT NOT NULL,
    product_id INT NOT NULL,
    CONSTRAINT fk_user_transactions FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_product_transactions FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE ON UPDATE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

COMMIT;