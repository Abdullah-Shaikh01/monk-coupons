-- Drop tables if they already exist
DROP TABLE IF EXISTS coupon_get_products;
DROP TABLE IF EXISTS coupon_buy_products;
DROP TABLE IF EXISTS coupons;
DROP TABLE IF EXISTS products;

-- Products table
CREATE TABLE IF NOT EXISTS products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL
);

-- Coupons table
CREATE TABLE IF NOT EXISTS coupons (
    id INT AUTO_INCREMENT PRIMARY KEY,
    type ENUM('cart-wise', 'product-wise', 'bxgy') NOT NULL,
    discount_value DECIMAL(10, 2),
    discount_type ENUM('percentage', 'fixed') NOT NULL DEFAULT 'percentage',
    quantity INT, -- Used in bxgy as required buy quantity
    repetition_threshold INT, -- Used as repetition_limit (bxgy) or min cart value (cart-wise)
    product_id INT, -- For product-wise coupons
    expiration_date DATE DEFAULT (CURRENT_DATE + INTERVAL 1 MONTH),
    FOREIGN KEY (product_id) REFERENCES products(id)
);

-- Buy products for bxgy coupons
CREATE TABLE IF NOT EXISTS coupon_buy_products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    coupon_id INT NOT NULL,
    product_id INT NOT NULL,
    FOREIGN KEY (coupon_id) REFERENCES coupons(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

-- Get products for bxgy coupons
CREATE TABLE IF NOT EXISTS coupon_get_products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    coupon_id INT NOT NULL,
    product_id INT NOT NULL,
    FOREIGN KEY (coupon_id) REFERENCES coupons(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

-- Seed: 10 products
INSERT INTO products (name, price) VALUES 
('Laptop', 75000.00),
('Smartphone', 30000.00),
('Headphones', 2500.00),
('Keyboard', 1500.00),
('Mouse', 800.00),
('Monitor', 12000.00),
('USB Cable', 300.00),
('Charger', 1000.00),
('Power Bank', 2500.00),
('Webcam', 3500.00);
