CREATE DATABASE IF NOT EXISTS NoPenNoPaper CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE NoPenNoPaper;

CREATE TABLE IF NOT EXISTS characters (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    profession VARCHAR(50) NOT NULL,
    age INTEGER NOT NULL,
    gender VARCHAR(10) NOT NULL,
    residence VARCHAR(50) NOT NULL,
    birthplace VARCHAR(50) NOT NULL
);