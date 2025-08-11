-- Minimal schema (生产建议用 GORM AutoMigrate 或专业迁移工具)
CREATE TABLE IF NOT EXISTS users (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  mobile VARCHAR(20) UNIQUE,
  password VARCHAR(255),
  nickname VARCHAR(50),
  avatar_url VARCHAR(255),
  created_at DATETIME,
  updated_at DATETIME
);

CREATE TABLE IF NOT EXISTS merchants (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(100),
  address VARCHAR(255),
  lat DOUBLE,
  lng DOUBLE,
  open_at VARCHAR(20),
  close_at VARCHAR(20),
  description VARCHAR(255),
  created_at DATETIME,
  updated_at DATETIME
);

CREATE TABLE IF NOT EXISTS dishes (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  merchant_id BIGINT UNSIGNED,
  category VARCHAR(50),
  name VARCHAR(100),
  price BIGINT,
  pic_url VARCHAR(255),
  stock INT,
  description VARCHAR(255),
  created_at DATETIME,
  updated_at DATETIME,
  INDEX idx_dish_merchant(merchant_id)
);

CREATE TABLE IF NOT EXISTS cart_items (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT UNSIGNED,
  merchant_id BIGINT UNSIGNED,
  dish_id BIGINT UNSIGNED,
  quantity INT,
  created_at DATETIME,
  updated_at DATETIME,
  INDEX idx_cart_user(user_id)
);

CREATE TABLE IF NOT EXISTS orders (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT UNSIGNED,
  merchant_id BIGINT UNSIGNED,
  status VARCHAR(20),
  total_amount BIGINT,
  pay_amount BIGINT,
  created_at DATETIME,
  updated_at DATETIME,
  INDEX idx_order_user(user_id)
);

CREATE TABLE IF NOT EXISTS order_items (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  order_id BIGINT UNSIGNED,
  dish_id BIGINT UNSIGNED,
  name VARCHAR(100),
  price BIGINT,
  quantity INT,
  created_at DATETIME,
  updated_at DATETIME,
  INDEX idx_oi_order(order_id)
);