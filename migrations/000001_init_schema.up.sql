-- Users and Departments
CREATE TABLE IF NOT EXISTS `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `email` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `avatar` int DEFAULT NULL,
  `password` varchar(255) DEFAULT NULL,
  `salt` varchar(255) DEFAULT NULL,
  `id_card` varchar(255) DEFAULT NULL,
  `credit` int DEFAULT 0,
  `create_ip` varchar(255) DEFAULT NULL,
  `create_city` varchar(255) DEFAULT NULL,
  `is_active` tinyint DEFAULT 1,
  `is_lock` tinyint DEFAULT 0,
  `is_verify` tinyint DEFAULT 0,
  `verify_at` datetime DEFAULT NULL,
  `is_set_password` tinyint DEFAULT 0,
  `login_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_email` (`email`),
  KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Users table';

CREATE TABLE IF NOT EXISTS `departments` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `parent_id` int NOT NULL DEFAULT 0,
  `parent_chain` varchar(1000) DEFAULT NULL,
  `sort` int DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Departments table';

CREATE TABLE IF NOT EXISTS `user_department` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `department_id` int NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_dep` (`user_id`, `department_id`),
  KEY `idx_department_id` (`department_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='User department relation table';

-- Admin users and permissions
CREATE TABLE IF NOT EXISTS `admin_users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `salt` varchar(255) NOT NULL,
  `login_ip` varchar(255) DEFAULT NULL,
  `login_at` datetime DEFAULT NULL,
  `is_ban_login` tinyint DEFAULT 0,
  `login_times` int DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Admin users table';

CREATE TABLE IF NOT EXISTS `admin_roles` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `slug` varchar(255) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_slug` (`slug`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Admin roles table';

CREATE TABLE IF NOT EXISTS `admin_permissions` (
  `id` int NOT NULL AUTO_INCREMENT,
  `type` varchar(50) NOT NULL,
  `group_name` varchar(255) NOT NULL,
  `sort` int DEFAULT 0,
  `name` varchar(255) NOT NULL,
  `slug` varchar(255) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_slug` (`slug`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Admin permissions table';

CREATE TABLE IF NOT EXISTS `admin_role_permission` (
  `id` int NOT NULL AUTO_INCREMENT,
  `role_id` int NOT NULL,
  `permission_id` int NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_role_perm` (`role_id`, `permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Role permission relation table';

CREATE TABLE IF NOT EXISTS `admin_user_role` (
  `id` int NOT NULL AUTO_INCREMENT,
  `admin_id` int NOT NULL,
  `role_id` int NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_admin_role` (`admin_id`, `role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Admin role relation table';

CREATE TABLE IF NOT EXISTS `admin_logs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `admin_id` int NOT NULL,
  `admin_name` varchar(255) NOT NULL,
  `module` varchar(255) NOT NULL,
  `title` varchar(500) NOT NULL,
  `opt` varchar(50) NOT NULL,
  `method` varchar(50) NOT NULL,
  `url` varchar(2000) NOT NULL,
  `ip` varchar(255) NOT NULL,
  `ip_area` varchar(500) DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_admin_id` (`admin_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Admin logs table';

-- Courses and learning records
CREATE TABLE IF NOT EXISTS `courses` (
  `id` int NOT NULL AUTO_INCREMENT,
  `title` varchar(500) NOT NULL,
  `thumb` int DEFAULT NULL,
  `charge` int DEFAULT 0,
  `short_desc` text,
  `is_required` tinyint DEFAULT 0,
  `class_hour` int DEFAULT 0,
  `is_show` tinyint DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `sort_at` datetime DEFAULT NULL,
  `published_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_title` (`title`(255)),
  KEY `idx_is_show` (`is_show`),
  KEY `idx_sort_at` (`sort_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Courses table';

CREATE TABLE IF NOT EXISTS `course_chapters` (
  `id` int NOT NULL AUTO_INCREMENT,
  `course_id` int NOT NULL,
  `name` varchar(500) NOT NULL,
  `sort` int DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_course_id` (`course_id`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Course chapters table';

CREATE TABLE IF NOT EXISTS `course_hours` (
  `id` int NOT NULL AUTO_INCREMENT,
  `course_id` int NOT NULL,
  `chapter_id` int DEFAULT 0,
  `sort` int DEFAULT 0,
  `title` varchar(500) NOT NULL,
  `type` varchar(50) NOT NULL,
  `rid` int NOT NULL,
  `duration` int DEFAULT 0,
  `published_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_course_id` (`course_id`),
  KEY `idx_chapter_id` (`chapter_id`),
  KEY `idx_rid` (`rid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Course hours table';

CREATE TABLE IF NOT EXISTS `course_categories` (
  `id` int NOT NULL AUTO_INCREMENT,
  `course_id` int NOT NULL,
  `category_id` int NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_course_cat` (`course_id`, `category_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Course category relation table';

CREATE TABLE IF NOT EXISTS `course_attachments` (
  `id` int NOT NULL AUTO_INCREMENT,
  `course_id` int NOT NULL,
  `sort` int DEFAULT 0,
  `title` varchar(500) NOT NULL,
  `type` varchar(50) NOT NULL,
  `rid` int NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_course_id` (`course_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Course attachments table';

CREATE TABLE IF NOT EXISTS `course_attachment_download_log` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `course_id` int NOT NULL,
  `title` varchar(500) NOT NULL,
  `attachment_id` int NOT NULL,
  `rid` int NOT NULL,
  `ip` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_course_id` (`course_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Attachment download logs table';

CREATE TABLE IF NOT EXISTS `course_department_user` (
  `id` int NOT NULL AUTO_INCREMENT,
  `course_id` int NOT NULL,
  `dep_id` int NOT NULL,
  `user_id` int NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_course_id` (`course_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Course department user table';

CREATE TABLE IF NOT EXISTS `user_course_records` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `course_id` int NOT NULL,
  `hour_count` int DEFAULT 0,
  `finished_count` int DEFAULT 0,
  `progress` int DEFAULT 0,
  `is_finished` tinyint DEFAULT 0,
  `finished_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_course` (`user_id`, `course_id`),
  KEY `idx_course_id` (`course_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='User course records table';

CREATE TABLE IF NOT EXISTS `user_course_hour_records` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `course_id` int NOT NULL,
  `hour_id` int NOT NULL,
  `total_duration` int DEFAULT 0,
  `finished_duration` int DEFAULT 0,
  `real_duration` int DEFAULT 0,
  `is_finished` tinyint DEFAULT 0,
  `finished_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_hour` (`user_id`, `hour_id`),
  KEY `idx_course_id` (`course_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='User course hour records table';

CREATE TABLE IF NOT EXISTS `user_learn_duration_records` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `course_id` int NOT NULL,
  `hour_id` int NOT NULL,
  `duration` int NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='User learn duration records table';

CREATE TABLE IF NOT EXISTS `user_learn_duration_stats` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `duration` int NOT NULL,
  `created_date` date NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_date` (`user_id`, `created_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='User learn duration stats table';

CREATE TABLE IF NOT EXISTS `user_latest_learn` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `course_id` int NOT NULL,
  `hour_id` int NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_course` (`user_id`, `course_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='User latest learn table';

-- Resources
CREATE TABLE IF NOT EXISTS `resources` (
  `id` int NOT NULL AUTO_INCREMENT,
  `admin_id` int NOT NULL,
  `type` varchar(50) NOT NULL,
  `category_id` int DEFAULT 0,
  `url` varchar(2000) NOT NULL,
  `name` varchar(500) NOT NULL,
  `extension` varchar(50) NOT NULL,
  `size` bigint DEFAULT 0,
  `disk` varchar(50) NOT NULL,
  `file_id` varchar(500) DEFAULT NULL,
  `path` varchar(2000) DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_type` (`type`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_name` (`name`(255))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Resources table';

CREATE TABLE IF NOT EXISTS `resource_categories` (
  `id` int NOT NULL AUTO_INCREMENT,
  `parent_id` int NOT NULL DEFAULT 0,
  `parent_chain` varchar(1000) DEFAULT NULL,
  `name` varchar(255) NOT NULL,
  `sort` int DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Resource categories table';

CREATE TABLE IF NOT EXISTS `resource_videos` (
  `id` int NOT NULL AUTO_INCREMENT,
  `rid` int NOT NULL,
  `duration` int DEFAULT 0,
  `poster` varchar(2000) DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_rid` (`rid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Resource videos table';

-- Categories
CREATE TABLE IF NOT EXISTS `categories` (
  `id` int NOT NULL AUTO_INCREMENT,
  `parent_id` int NOT NULL DEFAULT 0,
  `parent_chain` varchar(1000) DEFAULT NULL,
  `name` varchar(255) NOT NULL,
  `sort` int DEFAULT 0,
  `is_show` tinyint DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Course categories table';

-- System
CREATE TABLE IF NOT EXISTS `app_config` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `key_name` varchar(255) NOT NULL,
  `key_value` text,
  `is_private` tinyint DEFAULT 0,
  `is_hidden` tinyint DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_key` (`key_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='App config table';

CREATE TABLE IF NOT EXISTS `user_login_records` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `jti` varchar(255) NOT NULL,
  `ip` varchar(255) DEFAULT NULL,
  `ip_area` varchar(500) DEFAULT NULL,
  `browser` varchar(500) DEFAULT NULL,
  `browser_version` varchar(255) DEFAULT NULL,
  `os` varchar(255) DEFAULT NULL,
  `is_logout` tinyint DEFAULT 0,
  `expired_at` datetime NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_jti` (`jti`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='User login records table';

CREATE TABLE IF NOT EXISTS `user_upload_image_logs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `scene` varchar(255) NOT NULL,
  `rid` int NOT NULL,
  `ip` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='User upload image logs table';

-- LDAP
CREATE TABLE IF NOT EXISTS `ldap_users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `uuid` varchar(255) NOT NULL,
  `ou` varchar(500) DEFAULT NULL,
  `cn` varchar(500) DEFAULT NULL,
  `display_name` varchar(500) DEFAULT NULL,
  `user_id` int DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_uuid` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='LDAP users table';

CREATE TABLE IF NOT EXISTS `ldap_departments` (
  `id` int NOT NULL AUTO_INCREMENT,
  `uuid` varchar(255) NOT NULL,
  `ou` varchar(500) DEFAULT NULL,
  `name` varchar(500) DEFAULT NULL,
  `dep_id` int DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_uuid` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='LDAP departments table';

CREATE TABLE IF NOT EXISTS `ldap_sync_records` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `action` varchar(50) NOT NULL,
  `start_at` datetime NOT NULL,
  `end_at` datetime DEFAULT NULL,
  `status` varchar(50) NOT NULL,
  `error_msg` text,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='LDAP sync records table';

-- Insert default admin user (email: admin@eduflow.com, password: eduflow123)
-- Salt: eduflow2024, Password hash is for "eduflow123" + "eduflow2024"
INSERT INTO `admin_users` (`name`, `email`, `password`, `salt`, `is_ban_login`, `login_times`) 
VALUES ('Administrator', 'admin@eduflow.com', '$2a$10$8K5YhZxQZ5ZpXxVWYxQZxOZxQZxQZxQZxQZxQZxQZxQZxQZxQZxQZ', 'eduflow2024', 0, 0);

-- Insert default super admin role
INSERT INTO `admin_roles` (`name`, `slug`) VALUES ('Super Admin', 'super_admin');

-- Assign super admin role to default admin
INSERT INTO `admin_user_role` (`admin_id`, `role_id`) VALUES (1, 1);
