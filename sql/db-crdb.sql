-- This file describes the schemas used for the app for CockroachDB databases.

-- Create the sequence that will set the primary key of rows in the `sites` table.
CREATE SEQUENCE sites_id;

CREATE TABLE sites (
  id INT PRIMARY KEY DEFAULT nextval('sites_id'),
  domain STRING NOT NULL UNIQUE, -- Able to fit in a cache line along with the length.
  name STRING NOT NULL,
  logo STRING NOT NULL DEFAULT '',
  favicon STRING NOT NULL DEFAULT '',
  tls INT NOT NULL DEFAULT 0,
  created TIMESTAMP NOT NULL DEFAULT now()
);

CREATE SEQUENCE users_id;

CREATE TABLE users (
  id INT PRIMARY KEY DEFAULT nextval('users_id'),
  username STRING NOT NULL UNIQUE,
  email STRING NOT NULL UNIQUE,
  pass BYTES NOT NULL,
  fname STRING NOT NULL,
  lname STRING NOT NULL,
  updated TIMESTAMP NOT NULL DEFAULT now(), -- TODO: ON UPDATE CURRENT_TIMESTAMP
  registered TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE usermeta (
  user_id INT,
  k STRING,
  v STRING NOT NULL,
  updated TIMESTAMP NOT NULL DEFAULT now(), -- TODO: ON UPDATE CURRENT_TIMESTAMP
  PRIMARY KEY (user_id, k),
  CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE SEQUENCE content_id;

-- The kinds of things that will appear in the "content" table are things that are viewable from a URL; that is, they have a slug that will be part of some URL.
-- A row in this table may represent an item that only functions as a parent to other items, where the parent slug is used in the URL along with the child slug.
CREATE TABLE content (
  id INT PRIMARY KEY DEFAULT nextval('content_id'),
  site INT NOT NULL,
  slug STRING NOT NULL CHECK (length(slug) > 0),
  author INT NOT NULL,
  type STRING NOT NULL CHECK (length(type) > 0),
  parent INT NOT NULL DEFAULT 0,
  title STRING NOT NULL,
  meta_title STRING NOT NULL DEFAULT '',
  meta_desc STRING NOT NULL DEFAULT '',
  body BYTES NOT NULL DEFAULT '',
  status STRING NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'unsaved', 'trashed')) ,
  updated TIMESTAMP NOT NULL DEFAULT now(),
  CONSTRAINT fk_site_id FOREIGN KEY (site) REFERENCES sites (id),
  CONSTRAINT fk_author_id FOREIGN KEY (author) REFERENCES users (id),
  UNIQUE (site, slug),
  INDEX site_type_status (site, type, status),
  FAMILY f1 (id, site, slug, author, type, parent, title, meta_title, meta_desc),
  FAMILY f2 (body, status, updated)
);

CREATE SEQUENCE blobs_id;

-- Arbitrary data keyed by 'role' string and/or 'k' integer.
-- Such a table is replicated for each website.
CREATE TABLE blobs (
  id INT PRIMARY KEY DEFAULT nextval('blobs_id'),
  role STRING NOT NULL,
  k INT NOT NULL DEFAULT 0, -- same type as id of table content but not necessarily a foreign key
  v BYTES NOT NULL,
  updated TIMESTAMP NOT NULL DEFAULT now(),
  INDEX role_k (role, k)
);

-- The options table contains basic settings that can be set on any site or globally for all sites (with site = main_site_id).
CREATE TABLE options (
  site INT,
  k STRING,
  v BYTES NOT NULL,
  PRIMARY KEY (site, k),
  CONSTRAINT fk_site_id FOREIGN KEY (site) REFERENCES sites (id)
);

CREATE SEQUENCE site_messages_id;

-- The site_messages table contains alert messages pertaining to sites.
CREATE TABLE site_messages (
  id INT PRIMARY KEY DEFAULT nextval('site_messages_id'),
  site INT,
  role STRING NOT NULL,
  k STRING NOT NULL,
  v STRING NOT NULL,
  created TIMESTAMP NOT NULL DEFAULT now(),
  CONSTRAINT fk_site_id FOREIGN KEY (site) REFERENCES sites (id),
  INDEX indx_role (role),
  INDEX indx_k (k)
);

CREATE SEQUENCE user_messages_id;

-- The user_messages table contains alert messages pertaining to sites.
CREATE TABLE user_messages (
  id INT PRIMARY KEY DEFAULT nextval('user_messages_id'),
  user_id INT,
  k STRING,
  v STRING NOT NULL,
  created TIMESTAMP NOT NULL DEFAULT now(),
  CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);

-- END (THIS LINE IS NEEDED FOR THE SETUP TOOL TO STOP PARSING THE REST OF THE FILE)

CREATE TABLE media (
  id STRING PRIMARY KEY,
  ext STRING NOT NULL DEFAULT '', -- either .extension (including the dot) or blank
  site INT NOT NULL,
  name STRING NOT NULL DEFAULT '',
  alt STRING NOT NULL DEFAULT '',
  description STRING NOT NULL DEFAULT '',
  uploaded TIMESTAMP NOT NULL DEFAULT now(), -- TODO: ON UPDATE CURRENT_TIMESTAMP
  CONSTRAINT fk_site_id FOREIGN KEY (site) REFERENCES sites (id)
);

CREATE TABLE `products` (
  `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `name` STRING NOT NULL,
  `description` TEXT NOT NULL,
  `confirmation` TEXT NOT NULL,
  `initial_price` DECIMAL(9,2) NOT NULL DEFAULT 0.00,
  `recur_price` DECIMAL(9,2) NOT NULL DEFAULT 0.00,
  `cycle_number` INT NOT NULL DEFAULT 1,
  `cycle_type` ENUM('day','week','month','year') NOT NULL DEFAULT 'month',
  `trial_price` DECIMAL(9,2) NOT NULL,
  `trial_cycles` INT NOT NULL DEFAULT 0,
  `active` INT NOT NULL DEFAULT 1,
  `expiration_cycles_number` INT NOT NULL, -- how many billing cycles this product is valid for
  `created` TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE `orders` (
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `user_id` INT NOT NULL,
  `product_id` INT NOT NULL,  -- Foreign key: member_products.id
  `bill_name` STRING NOT NULL,
  `bill_street` STRING NOT NULL,
  `bill_city` STRING NOT NULL,
  `bill_state` STRING NOT NULL,
  `bill_zip` STRING NOT NULL,
  `bill_country` STRING NOT NULL,
  `payment_type` STRING NOT NULL,
  `card_type` STRING NOT NULL,
  `last_four` STRING NOT NULL,
  `expire_month` CHAR(2) NOT NULL,
  `expire_year` STRING NOT NULL,
  `status` STRING NOT NULL,
  `transaction_id` STRING NOT NULL,
  `subtotal` DECIMAL(9,2) NOT NULL DEFAULT 0.00,
  `total` DECIMAL(9,2) NOT NULL DEFAULT 0.00,
  `discount_code` INT NOT NULL,  -- Foreign key: member_discount_codes.id
  `created` TIMESTAMP NOT NULL DEFAULT now(),
  FOREIGN KEY (`discount_code`) REFERENCES (`discount_codes`)
);

-- All of the registered discount codes.
CREATE TABLE `discount_codes` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `code` STRING NOT NULL,
  `starts` DATE NOT NULL,
  `expires` DATE NOT NULL,
  `max_uses` INT NOT NULL DEFAULT -1, -- unlimited uses allowed by default
  KEY (`code`)
);

-- The products to which each of the discount codes applies.
-- These records override any values set in the `products` table.
CREATE TABLE `discount_code_products` (
  `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `code_id` INT NOT NULL, -- Foreign key: member_discount_codes.id
  `product_id` INT NOT NULL,  -- Foreign key: member_products.id
  `initial_price` DECIMAL(9,2) NOT NULL DEFAULT 0.00,
  `recur_price` DECIMAL(9,2) NOT NULL DEFAULT 0.00,
  `cycle_number` INT NOT NULL,
  `cycle_type` ENUM('day','week','month','year') NOT NULL,
  `trial_price` DECIMAL(9,2) NOT NULL,
  `trial_cycles` INT NOT NULL,
  `expiration_cycles_number` INT NOT NULL,
  `created` TIMESTAMP NOT NULL DEFAULT now(),
  FOREIGN KEY `code_id` REFERENCES `discount_codes` (`product`)
);
