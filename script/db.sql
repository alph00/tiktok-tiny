-- username:gorm  
-- password:gorm  
-- addr:127.0.0.1:3306


CREATE DATABASE  IF NOT EXISTS `gorm` 
USE `gorm`;

-- message
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

-- user
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_name` varchar(40) NOT NULL,
  `password` varchar(256) NOT NULL,
  `following_count` bigint unsigned NOT NULL DEFAULT '0',
  `follower_count` bigint unsigned NOT NULL DEFAULT '0',
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `avatar` varchar(256) DEFAULT NULL,
  `background_image` varchar(256) DEFAULT 'default_background.jpg',
  `work_count` bigint unsigned NOT NULL DEFAULT '0',
  `favorite_count` bigint unsigned NOT NULL DEFAULT '0',
  `total_favorited` bigint unsigned NOT NULL DEFAULT '0',
  `signature` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `users_user_name_uindex` (`user_name`),
  KEY `idx_username` (`user_name`),
  KEY `idx_user_deleted_at` (`deleted_at`),
  KEY `idx_users_deleted_at` (`deleted_at`)
);