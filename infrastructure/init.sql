CREATE DATABASE IF NOT EXISTS unicorndb;
USE unicorndb;

CREATE TABLE IF NOT EXISTS unicorns (
    id INT AUTO_INCREMENT PRIMARY KEY,
    redis_token varchar(255) NOT NULL,
    filesystem_token varchar(255) NOT NULL
);

CREATE USER IF NOT EXISTS 'username'@'%' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON unicorndb.* to 'username'@'%';
FLUSH PRIVILEGES;
