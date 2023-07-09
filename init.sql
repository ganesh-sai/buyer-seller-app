
CREATE TABLE IF NOT EXISTS sellers (
                         id INT PRIMARY KEY AUTO_INCREMENT,
                         name VARCHAR(255) NOT NULL,
                         location VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS  products (
                          id INT PRIMARY KEY AUTO_INCREMENT,
                          seller_id INT NOT NULL,
                          product_name VARCHAR(255) NOT NULL,
                          price DECIMAL(10, 2) NOT NULL,
                          quantity INT NOT NULL,
                          CONSTRAINT fk_seller_id FOREIGN KEY (seller_id) REFERENCES sellers(id)
);

CREATE INDEX idx_product_seller_id ON products (seller_id);

SELECT * FROM db.products limit 1;