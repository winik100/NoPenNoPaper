USE NoPenNoPaper;

SOURCE internal/sql/users.sql;
SOURCE internal/sql/characters.sql;

INSERT INTO users (name, hashed_password, role) VALUES ('nadine', '$2a$12$LwKInuN4IE0pJqhGQWWUgespgEPI0302Hnts88fzAZZs54yflbJXO', 'gm');