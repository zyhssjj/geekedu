-- 创建数据库（不存在则创建）
CREATE DATABASE IF NOT EXISTS geekedu DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 切换到geekedu数据库
USE geekedu;

-- 新增：创建用户表（若不存在），role默认0=学生
CREATE TABLE IF NOT EXISTS `users` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(100) NOT NULL COMMENT 'bcrypt加密后的密码',
  `role` tinyint NOT NULL DEFAULT 0 COMMENT '角色：0-学生（默认），1-管理员',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 插入管理员用户（用户名：admin，密码：admin123，已bcrypt加密，满足快速登录需求）
INSERT INTO `users` (`username`, `password`, `role`, `created_at`, `updated_at`)
VALUES (
  'admin',
  '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- bcrypt加密后的admin123
  1,
  NOW(),
  NOW()
) ON DUPLICATE KEY UPDATE `updated_at` = NOW();