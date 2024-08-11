CREATE DATABASE IF NOT EXISTS NoPenNoPaper CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE NoPenNoPaper;

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

CREATE TABLE IF NOT EXISTS sessions (
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);

-- users.sql
CREATE TABLE IF NOT EXISTS users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(30) NOT NULL,
    hashed_password VARCHAR(60) NOT NULL,
    role VARCHAR(10) NOT NULL
);

INSERT INTO users (name, hashed_password, role) VALUES ('testgm', '$2a$12$EeAcZSu5HYgydNKQVKaAW.qdBMSNVEeGugDA1yoyrMQF12BZTxf76', 'gm');

CREATE TABLE IF NOT EXISTS materials (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    file_name VARCHAR(100) NOT NULL,
    uploaded_by INTEGER NOT NULL,
    CONSTRAINT fk_users_materials_id FOREIGN KEY (uploaded_by) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT unique_filename_user UNIQUE (file_name, uploaded_by)
);

-- characters.sql
CREATE TABLE IF NOT EXISTS characters (
	id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	created_by INTEGER NOT NULL,
	FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS character_info (
	character_id INTEGER NOT NULL PRIMARY KEY,
	name VARCHAR(50) NOT NULL,
	profession VARCHAR(50) NOT NULL,
	age INTEGER NOT NULL,
	gender VARCHAR(10) NOT NULL,
	residence VARCHAR(50) NOT NULL,
	birthplace VARCHAR(50) NOT NULL,
	FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE
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
	FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS character_stats (
	character_id INTEGER NOT NULL PRIMARY KEY,
	maxtp INTEGER NOT NULL,
	tp INTEGER NOT NULL,
	maxsta INTEGER NOT NULL,
	sta INTEGER NOT NULL,
	maxmp INTEGER NOT NULL,
	mp INTEGER NOT NULL,
	maxluck INTEGER NOT NULL,
	luck INTEGER NOT NULL,
	FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE
); 

CREATE TABLE IF NOT EXISTS skills (
	name VARCHAR(50) NOT NULL PRIMARY KEY,
	default_value INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS character_skills (
	character_id INTEGER NOT NULL,
	skill_name VARCHAR(50) NOT NULL,
	value INTEGER NOT NULL,
	CONSTRAINT fk_character_cs FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE,
	CONSTRAINT fk_skill_cs FOREIGN KEY (skill_name) REFERENCES skills(name),
	CONSTRAINT pk_character_skills PRIMARY KEY (character_id, skill_name)
);

CREATE TABLE IF NOT EXISTS custom_skills (
	name VARCHAR(50) NOT NULL PRIMARY KEY,
	category VARCHAR(50) NOT NULL,
	default_value INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS character_custom_skills (
	character_id INTEGER NOT NULL,
	custom_skill_name VARCHAR(50) NOT NULL,
	value INTEGER NOT NULL,
	CONSTRAINT fk_character_ccs FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE,
	CONSTRAINT fk_custom_skill_ccs FOREIGN KEY (custom_skill_name) REFERENCES custom_skills(name),
	CONSTRAINT pk_character_custom_skills PRIMARY KEY (character_id, custom_skill_name)
);


CREATE TABLE IF NOT EXISTS items (
	item_id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	character_id INTEGER NOT NULL,
	name VARCHAR(50) NOT NULL,
	description VARCHAR(255) NOT NULL,
	cnt INTEGER NOT NULL,
	FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS notes (
	note_id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	character_id INTEGER NOT NULL,
	text VARCHAR(255) NOT NULL,
	FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE
);

-- populate.sql
INSERT INTO skills (name, default_value) VALUES ('Anthropologie', 1),
			('Archäologie', 1),
			('Autofahren', 20),
			('Bibliotheksnutzung', 20),
			('Buchführung', 5),
			('Charme', 15),
			('Cthulhu-Mythos', 0),
			('Einschüchtern', 15),
			('Elektrische Reparaturen', 10),
			('Erste Hilfe', 30),
			('Finanzkraft', 0),
			('Geschichte', 5),
			('Horchen', 20),
			('Kaschieren', 10),
			('Klettern', 20),
			('Mechanische Reparaturen', 10),
			('Medizin', 1),
			('Naturkunde', 10),
			('Okkultismus', 5),
			('Orientierung', 10),
			('Psychoanalyse', 1),
			('Psychologie', 10),
			('Rechtswesen', 5),
			('Reiten', 5),
			('Schließtechnik', 1),
			('Schweres Gerät', 1),
			('Schwimmen', 20),
			('Springen', 20),
			('Spurensuche', 10),
			('Überreden', 5),
			('Überzeugen', 10),
			('Verborgen bleiben', 20),
			('Verborgenes erkennen', 25),
			('Verkleiden', 5),
			('Werfen', 20),
			('Werte schätzen', 5);