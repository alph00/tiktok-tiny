-- username:gorm  
-- password:gorm  
-- addr:127.0.0.1:3306


CREATE DATABASE  IF NOT EXISTS `gorm` 
USE `gorm`;

DROP TABLE IF EXISTS `messages`;
CREATE TABLE `messages` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `from_user_id` bigint unsigned NOT NULL,
  `to_user_id` bigint unsigned NOT NULL,
  `content` varchar(255) NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_userid_from` (`from_user_id`),
  KEY `idx_userid_to` (`to_user_id`),
  KEY `idx_messages_deleted_at` (`deleted_at`),
  KEY `idx_messages_created_at` (`created_at`)); 