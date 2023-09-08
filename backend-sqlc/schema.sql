CREATE TABLE `roles` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `status` enum('DISABLED','ENABLED','DRAFT') NOT NULL DEFAULT 'ENABLED',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `last_updated_by` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE permissions (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `role_id` bigint NOT NULL,
  `module_name` varchar(255) NOT NULL,
  `create` boolean NOT NULL DEFAULT 0,
  `read` boolean NOT NULL DEFAULT 0,
  `update` boolean NOT NULL DEFAULT 0,
  `delete` boolean NOT NULL DEFAULT 0,
  `export` boolean NOT NULL DEFAULT 0,
  `import` boolean NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_role_permissions` FOREIGN KEY (`role_id`) REFERENCES `roles`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `users` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `username` varchar(255) DEFAULT NULL,
  `password` varchar(255) NOT NULL,
  `staff_id` varchar(255) DEFAULT NULL,
  `email` varchar(255) NOT NULL,
  `email_confirmation` tinyint(1) DEFAULT '0',
  `status` enum('DISABLED','ENABLED','DRAFT') NOT NULL DEFAULT 'ENABLED',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `last_updated_by` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `roles_users` (
  `role_id` bigint NOT NULL,
  `user_id` bigint NOT NULL,
  PRIMARY KEY (`role_id`, `user_id`),
  CONSTRAINT `fk_user_roles_roles` FOREIGN KEY (`role_id`) REFERENCES `roles`(`id`),
  CONSTRAINT `fk_user_roles_users` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
