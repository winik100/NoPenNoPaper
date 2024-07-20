CREATE USER 'web'@'localhost';
GRANT SELECT, INSERT, UPDATE, DELETE ON NoPenNoPaper.* TO 'web'@'localhost';
ALTER USER 'web'@'localhost' IDENTIFIED BY 'mellon';

INSERT INTO users (name, hashed_password, role) VALUES ('nadine', '$2a$12$LwKInuN4IE0pJqhGQWWUgespgEPI0302Hnts88fzAZZs54yflbJXO', 'gm');