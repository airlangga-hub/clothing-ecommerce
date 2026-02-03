CREATE TABLE `products` (
  `id` INTEGER AUTO_INCREMENT PRIMARY KEY,
  `name` VARCHAR(100) NOT NULL,
  `description` VARCHAR(100) NOT NULL,
  `price` INTEGER NOT NULL
);

CREATE TABLE `users` (
  `id` INTEGER AUTO_INCREMENT PRIMARY KEY,
  `email` VARCHAR(100) NOT NULL UNIQUE,
  `password` VARCHAR(100) NOT NULL,
  `role` VARCHAR(100) NOT NULL DEFAULT 'user'
);

CREATE TABLE `cart_items` (
  `id` INTEGER AUTO_INCREMENT PRIMARY KEY,
  `user_id` INTEGER NOT NULL,
  `product_id` INTEGER NOT NULL,
  `quantity` INTEGER NOT NULL DEFAULT 1,
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  FOREIGN KEY (`product_id`) REFERENCES `products` (`id`),
  UNIQUE (`user_id`, `product_id`)
);

CREATE TABLE `order_items` (
  `id` INTEGER AUTO_INCREMENT PRIMARY KEY,
  `order_id` INTEGER NOT NULL,
  `product_id` INTEGER NOT NULL,
  `quantity` INTEGER NOT NULL,
  FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`),
  FOREIGN KEY (`product_id`) REFERENCES `products` (`id`),
  UNIQUE (`order_id`, `product_id`)
);

CREATE TABLE `orders` (
  `id` INTEGER AUTO_INCREMENT PRIMARY KEY,
  `user_id` INTEGER NOT NULL,
  `total_price` INTEGER NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
);
