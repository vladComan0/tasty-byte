CREATE SCHEMA IF NOT EXISTS `tastybyte`;
USE `tastybyte`;

GRANT CREATE, SELECT, INSERT, UPDATE, DELETE, REFERENCES ON `tastybyte`.* TO `tastybyte_user`@`%`;

CREATE TABLE `recipes` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `description` text NOT NULL,
  `instructions` text NOT NULL,
  `preparation_time` varchar(10) NOT NULL,
  `cooking_time` varchar(10) NOT NULL,
  `portions` int NOT NULL,
  `created` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `tags` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `tag_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `recipe_tags` (
  `recipe_id` int NOT NULL,
  `tag_id` int NOT NULL,
  PRIMARY KEY (`recipe_id`, `tag_id`),
  FOREIGN KEY (`recipe_id`) REFERENCES `recipes`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`tag_id`) REFERENCES `tags`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- CREATE TABLE `sessions` (
--   `token` char(43) COLLATE utf8mb4_unicode_ci NOT NULL,
--   `data` blob NOT NULL,
--   `expiry` timestamp(6) NOT NULL,
--   PRIMARY KEY (`token`),
--   KEY `sessions_expiry_idx` (`expiry`)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- CREATE TABLE `users` (
--   `id` int NOT NULL AUTO_INCREMENT,
--   `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
--   `email` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
--   `hashed_password` char(60) COLLATE utf8mb4_unicode_ci NOT NULL,
--   `created` datetime NOT NULL,
--   PRIMARY KEY (`id`),
--   UNIQUE KEY `user_uc_email` (`email`)
-- ) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
