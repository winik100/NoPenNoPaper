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
	tp INTEGER NOT NULL,
	sta INTEGER NOT NULL,
	mp INTEGER NOT NULL,
	luck INTEGER NOT NULL,
	FOREIGN KEY (character_id) REFERENCES characters(id)

); 
CREATE TABLE IF NOT EXISTS character_skills (
	character_id INTEGER NOT NULL PRIMARY KEY,
	anthropology INTEGER NOT NULL,
	archaeology INTEGER NOT NULL,
	driving INTEGER NOT NULL,
	libraryResearch INTEGER NOT NULL,
	accounting INTEGER NOT NULL,
	charme INTEGER NOT NULL,
	cthulhuMythos INTEGER NOT NULL,
	intimidate INTEGER NOT NULL,
	electricRepairs INTEGER NOT NULL,
	firstAid INTEGER NOT NULL,
	financials INTEGER NOT NULL,
	history INTEGER NOT NULL,
	listening INTEGER NOT NULL,
	concealing INTEGER NOT NULL,
	climbing INTEGER NOT NULL,
	mechanicalRepairs INTEGER NOT NULL,
	medicine INTEGER NOT NULL,
	naturalHistory INTEGER NOT NULL,
	occultism INTEGER NOT NULL,
	orientation INTEGER NOT NULL,
	psychoAnalysis INTEGER NOT NULL,
	psychology INTEGER NOT NULL,
	law INTEGER NOT NULL,
	horseriding INTEGER NOT NULL,
	locks INTEGER NOT NULL,
	heavyMachinery INTEGER NOT NULL,
	swimming INTEGER NOT NULL,
	jumping INTEGER NOT NULL,
	tracking INTEGER NOT NULL,
	persuasion INTEGER NOT NULL,
	convincing INTEGER NOT NULL,
	stealth INTEGER NOT NULL,
	detectingSecrets INTEGER NOT NULL,
	disguising INTEGER NOT NULL,
	throwing INTEGER NOT NULL,
	valuation INTEGER NOT NULL,
	FOREIGN KEY (character_id) REFERENCES characters(id)
);
CREATE TABLE IF NOT EXISTS items (
	item_id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	character_id INTEGER NOT NULL,
	name VARCHAR(50) NOT NULL,
	description VARCHAR(255) NOT NULL,
	cnt INTEGER NOT NULL,
	FOREIGN KEY (character_id) REFERENCES characters(id)
)