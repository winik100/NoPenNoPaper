DELETE FROM character_attributes;
DELETE FROM character_info;
DELETE FROM character_stats;
DELETE FROM characters;
DELETE FROM sessions;
ALTER TABLE characters AUTO_INCREMENT=1;

DELETE FROM users;
ALTER TABLE users AUTO_INCREMENT=1;

DELETE FROM sessions;