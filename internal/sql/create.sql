CREATE DATABASE IF NOT EXISTS NoPenNoPaper CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE NoPenNoPaper;

CREATE TABLE IF NOT EXISTS characters (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT
);

CREATE TABLE IF NOT EXISTS character_info (
    character_id INTEGER NOT NULL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    profession VARCHAR(50) NOT NULL,
    age INTEGER NOT NULL,
    gender VARCHAR(10) NOT NULL,
    residence VARCHAR(50) NOT NULL,
    birthplace VARCHAR(50) NOT NULL,
    FOREIGN KEY (character_id) REFERENCES characters(id)
);

CREATE TABLE IF NOT EXISTS character_attributes (
    character_id INTEGER NOT NULL PRIMARY KEY,
    st INTEGER NOT NULL,
    ge INTEGER NOT NULL,
    ma INTEGER NOT NULL,
    ko INTEGER NOT NULL,
    er INTEGER NOT NULL,
    bi INTEGER NOT NULL,
    gr INTEGER NOT NULL,
    i INTEGER NOT NULL,
    bw INTEGER NOT NULL,
    FOREIGN KEY (character_id) REFERENCES characters(id)
);