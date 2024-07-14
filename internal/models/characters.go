package models

import (
	"database/sql"
	"errors"
)

type PersonalInfo struct {
	Name       string
	Profession string
	Age        string
	Gender     string
	Residence  string
	Birthplace string
}

type Attributes struct {
	ST int
	GE int
	MA int
	KO int
	ER int
	BI int
	GR int
	IN int
	BW int
}

type Character struct {
	ID         int
	Info       map[string]string
	Attributes map[string]int
}

type CharacterModelInterface interface {
	Insert(info map[string]string, attributes map[string]int) (int, error)
	Get(id int) (Character, error)
}

type CharacterModel struct {
	DB *sql.DB
}

func (c *CharacterModel) Insert(info map[string]string, attributes map[string]int) (int, error) {
	stmt := "INSERT INTO characters (id) VALUES (null);"
	result, err := c.DB.Exec(stmt)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	stmt = "INSERT INTO character_info (character_id, name, profession, age, gender, residence, birthplace) VALUES (?,?,?,?,?,?,?);"
	_, err = c.DB.Exec(stmt, id, info["name"], info["profession"], info["age"], info["gender"], info["residence"], info["birthplace"])
	if err != nil {
		return 0, err
	}

	stmt = "INSERT INTO character_attributes (character_id, st, ge, ma, ko, er, bi, gr, i, bw) VALUES (?,?,?,?,?,?,?,?,?,?);"
	_, err = c.DB.Exec(stmt, id, attributes["ST"], attributes["GE"], attributes["MA"], attributes["KO"], attributes["ER"], attributes["BI"], attributes["GR"], attributes["IN"], attributes["BW"])
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (c *CharacterModel) Get(id int) (Character, error) {
	var character Character

	stmt := "SELECT name, profession, age, gender, residence, birthplace FROM character_info WHERE character_id=?;"
	result := c.DB.QueryRow(stmt, id)
	err := result.Scan(character.Info["name"], character.Info["profession"], character.Info["age"], character.Info["gender"], character.Info["residence"], character.Info["birthplace"])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Character{}, ErrNoRecord
		}
		return Character{}, err
	}

	stmt = "SELECT st, ge, ma, ko, er, bi, gr, i, bw FROM character_attributes WHERE character_id=?;"
	result = c.DB.QueryRow(stmt, id)
	err = result.Scan(character.Attributes["ST"], character.Attributes["GE"], character.Attributes["MA"],
		character.Attributes["KO"], character.Attributes["ER"], character.Attributes["BI"],
		character.Attributes["GR"], character.Attributes["IN"], character.Attributes["BW"])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Character{}, ErrNoRecord
		}
		return Character{}, err
	}

	character.ID = id
	return character, nil
}
