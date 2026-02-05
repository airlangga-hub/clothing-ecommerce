CREATE TABLE products (
  id INTEGER AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  description VARCHAR(100) NOT NULL,
  price INTEGER NOT NULL,
  stock INTEGER NOT NULL
);

CREATE TABLE users (
  id INTEGER AUTO_INCREMENT PRIMARY KEY,
  email VARCHAR(100) NOT NULL UNIQUE,
  password VARCHAR(100) NOT NULL,
  role VARCHAR(50) NOT NULL DEFAULT 'user'
);

CREATE TABLE cart_items (
  id INTEGER AUTO_INCREMENT PRIMARY KEY,
  user_id INTEGER NOT NULL,
  product_id INTEGER NOT NULL,
  quantity INTEGER NOT NULL DEFAULT 1,
  FOREIGN KEY (user_id) REFERENCES users (id),
  FOREIGN KEY (product_id) REFERENCES products (id),
  UNIQUE (user_id, product_id)
);

CREATE TABLE orders (
  id INTEGER AUTO_INCREMENT PRIMARY KEY,
  user_id INTEGER NOT NULL,
  total_price INTEGER NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE order_items (
  id INTEGER AUTO_INCREMENT PRIMARY KEY,
  order_id INTEGER NOT NULL,
  product_id INTEGER NOT NULL,
  quantity INTEGER NOT NULL,
  FOREIGN KEY (order_id) REFERENCES orders (id),
  FOREIGN KEY (product_id) REFERENCES products (id),
  UNIQUE (order_id, product_id)
);

DELIMITER //
CREATE PROCEDURE place_order_items(IN o_id INTEGER, IN p_id INTEGER, IN qty INTEGER)
BEGIN
    DECLARE current_stock INTEGER;
    
    SELECT stock INTO current_stock FROM products WHERE id = p_id;
    
    IF current_stock < qty THEN
        SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Insufficient stock';
    ELSE
        INSERT INTO order_items
    		(order_id, product_id, quantity)
    	VALUES
    		(o_id, p_id, qty);
        
        UPDATE products
    	SET stock = stock - qty
    	WHERE id = p_id;
    END IF;
END //
DELIMITER ;