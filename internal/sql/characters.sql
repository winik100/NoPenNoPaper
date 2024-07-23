USE NoPenNoPaper;

CREATE TABLE IF NOT EXISTS characters (
	id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	created_by INTEGER NOT NULL,
	FOREIGN KEY (created_by) REFERENCES users(id)
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
	FOREIGN KEY (character_id) REFERENCES characters(id)
); 

CREATE TABLE IF NOT EXISTS skills (
	name VARCHAR(50) NOT NULL PRIMARY KEY,
	default_value INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS character_skills (
	character_id INTEGER NOT NULL,
	skill_name VARCHAR(50) NOT NULL,
	value INTEGER NOT NULL,
	CONSTRAINT fk_character FOREIGN KEY (character_id) REFERENCES characters(id),
	CONSTRAINT fk_skill FOREIGN KEY (skill_name) REFERENCES skills(name),
	CONSTRAINT character_skills_unique UNIQUE (character_id, skill_name)
);


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

CREATE TABLE IF NOT EXISTS items (
	item_id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	character_id INTEGER NOT NULL,
	name VARCHAR(50) NOT NULL,
	description VARCHAR(255) NOT NULL,
	cnt INTEGER NOT NULL,
	FOREIGN KEY (character_id) REFERENCES characters(id)
);

CREATE TABLE IF NOT EXISTS notes (
	character_id INTEGER NOT NULL,
	text VARCHAR(255) NOT NULL,
	FOREIGN KEY (character_id) REFERENCES characters(id)
);