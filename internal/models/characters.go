package models

import (
	"database/sql"
	"errors"
)

type PersonalInfo struct {
	Name       string
	Profession string
	Age        int
	Gender     string
	Residence  string
	Birthplace string
}

type Character struct {
	ID         int
	Info       PersonalInfo
	Attributes map[string]int
}

type CharacterModelInterface interface {
	Insert(info PersonalInfo, attributes map[string]int) (int, error)
	Get(id int) (Character, error)
	Latest() ([]Character, error)
}

type CharacterModel struct {
	DB *sql.DB
}

func (c *CharacterModel) Insert(info PersonalInfo, attributes map[string]int) (int, error) {
	stmt1 := "INSERT INTO characters (id) VALUES (null);"
	result, err := c.DB.Exec(stmt1)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	stmt2 := "INSERT INTO character_info (character_id, name, profession, age, gender, residence, birthplace) VALUES (?,?,?,?,?,?,?);"

	_, err = c.DB.Exec(stmt2, id, info.Name, info.Profession, info.Age, info.Gender, info.Residence, info.Birthplace)
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (c *CharacterModel) Get(id int) (Character, error) {
	stmt := "SELECT name, profession, age, gender, residence, birthplace FROM character_info WHERE id=?;"

	result := c.DB.QueryRow(stmt, id)
	var character Character

	err := result.Scan(&character.Info.Name, &character.Info.Profession, &character.Info.Age, &character.Info.Gender, &character.Info.Residence, &character.Info.Birthplace)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Character{}, ErrNoRecord
		}
		return Character{}, err
	}

	character.ID = id
	return character, nil
}

func (c *CharacterModel) Latest() ([]Character, error) {
	stmt := "SELECT character_id, name, profession, age, gender, residence, birthplace FROM character_info ORDER BY character_id DESC LIMIT 5;"

	rows, err := c.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var characters []Character
	for rows.Next() {
		var character Character
		err := rows.Scan(&character.ID, &character.Info.Name, &character.Info.Profession, &character.Info.Age, &character.Info.Gender, &character.Info.Residence, &character.Info.Birthplace)
		if err != nil {
			return nil, err
		}
		characters = append(characters, character)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return characters, nil
}
