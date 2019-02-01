-- 创建数据库
CREATE DATABASE IF NOT EXISTS `reptile`;

-- 使用数据库

USE `reptile`;

-- 创建表
CREATE TABLE IF NOT EXISTS `user_use`(
        `user_use_id` INT UNSIGNED AUTO_INCREMENT,
        `user_use_key` VARCHAR(100) NOT NULL,
        `user_use_count` VARCHAR(40) NOT NULL,
        `user_use_date` DATE,
        PRIMARY KEY ( `user_use_id` )
     )ENGINE=InnoDB DEFAULT CHARSET=utf8;
