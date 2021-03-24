-- sudo mysql < ./pkg/models/mysql/init.sql
drop database IF EXISTS shaderbox;

CREATE DATABASE shaderbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE shaderbox;

CREATE TABLE artworks (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    created DATETIME NOT NULL,
    rendering BOOL DEFAULT false,
    rendered BOOL DEFAULT false
);
-- Add an index on the created column.
CREATE INDEX idx_artworks_created ON artworks(created);


DROP USER IF EXISTS 'web'@'localhost';
FLUSH PRIVILEGES;
CREATE USER 'web'@'localhost';
GRANT SELECT, INSERT, UPDATE ON shaderbox.* TO 'web'@'localhost';
ALTER USER 'web'@'localhost' IDENTIFIED BY 'pass';

SHOW TABLES;
SELECT user from mysql.user;
