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
	CONSTRAINT fk_character_cs FOREIGN KEY (character_id) REFERENCES characters(id),
	CONSTRAINT fk_skill_cs FOREIGN KEY (skill_name) REFERENCES skills(name),
	CONSTRAINT pk_character_skills PRIMARY KEY (character_id, skill_name)
);

CREATE TABLE IF NOT EXISTS custom_skills (
	name VARCHAR(50) NOT NULL,
	category VARCHAR(50) NOT NULL,
	default_value INTEGER NOT NULL,
	CONSTRAINT pk_custom_skills PRIMARY KEY (name, category)
);

CREATE TABLE IF NOT EXISTS character_custom_skills (
	character_id INTEGER NOT NULL,
	custom_skill_name VARCHAR(50) NOT NULL,
	value INTEGER NOT NULL,
	CONSTRAINT fk_character_ccs FOREIGN KEY (character_id) REFERENCES characters(id),
	CONSTRAINT fk_custom_skill_ccs FOREIGN KEY (custom_skill_name) REFERENCES custom_skills(name),
	CONSTRAINT pk_character_custom_skills PRIMARY KEY (character_id, custom_skill_name)
);


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