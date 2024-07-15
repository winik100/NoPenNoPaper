CREATE DATABASE IF NOT EXISTS NoPenNoPaper CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE NoPenNoPaper;

CREATE TABLE IF NOT EXISTS sessions (
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);