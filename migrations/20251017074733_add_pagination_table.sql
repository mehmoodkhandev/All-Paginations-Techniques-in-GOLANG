-- Create "products" table
CREATE TABLE `products` (
  `id` char(36) NOT NULL DEFAULT (uuid()),
  `title` longtext NOT NULL,
  `description` longtext NOT NULL,
  `image` longtext NOT NULL,
  `price` bigint NOT NULL,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_products_deleted_at` (`deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
