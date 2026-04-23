CREATE DATABASE IF NOT EXISTS database;
USE database;

CREATE TABLE IF NOT EXISTS unicorns (
    id INT AUTO_INCREMENT PRIMARY KEY,
    redis_token varchar(255) NOT NULL,
    filesystem_token varchar(255) NOT NULL
);
