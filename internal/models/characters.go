package models

import (
	"database/sql"
	"errors"

	"github.com/justinian/dice"
)

type Character struct {
	ID                int
	Info              map[string]string
	Attributes        map[string]int
	DerivedAttributes map[string]int
}

type CharacterModelInterface interface {
	Insert(character Character) (int, error)
	Get(id int) (Character, error)
}

type CharacterModel struct {
	DB *sql.DB
}

func (c *CharacterModel) Insert(character Character) (int, error) {
	tx, err := c.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	stmt := "INSERT INTO characters (id) VALUES (null);"
	result, err := tx.Exec(stmt)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	stmt = "INSERT INTO character_info (character_id, name, profession, age, gender, residence, birthplace) VALUES (?,?,?,?,?,?,?);"
	_, err = tx.Exec(stmt, id, character.Info["name"], character.Info["profession"], character.Info["age"], character.Info["gender"], character.Info["residence"], character.Info["birthplace"])
	if err != nil {
		return 0, err
	}

	stmt = "INSERT INTO character_attributes (character_id, st, ge, ma, ko, er, bi, gr, i, bw) VALUES (?,?,?,?,?,?,?,?,?,?);"
	_, err = tx.Exec(stmt, id, character.Attributes["st"], character.Attributes["ge"], character.Attributes["ma"],
		character.Attributes["ko"], character.Attributes["er"], character.Attributes["bi"],
		character.Attributes["gr"], character.Attributes["in"], character.Attributes["bw"])
	if err != nil {
		return 0, err
	}

	derivedAttributes := character.deriveAttributes()
	stmt = "INSERT INTO character_stats (character_id, tp, sta, mp, luck) VALUES (?,?,?,?,?);"
	_, err = tx.Exec(stmt, id, derivedAttributes["tp"], derivedAttributes["sta"], derivedAttributes["mp"], derivedAttributes["luck"])
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
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
	err = result.Scan(character.Attributes["st"], character.Attributes["ge"], character.Attributes["ma"],
		character.Attributes["ko"], character.Attributes["er"], character.Attributes["bi"],
		character.Attributes["gr"], character.Attributes["in"], character.Attributes["bw"])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Character{}, ErrNoRecord
		}
		return Character{}, err
	}

	stmt = "SELECT tp, sta, mp, luck FROM character_stats WHERE character_id=?;"
	result = c.DB.QueryRow(stmt, id)
	err = result.Scan(character.DerivedAttributes["tp"], character.DerivedAttributes["sta"], character.DerivedAttributes["mp"], character.DerivedAttributes["luck"])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Character{}, ErrNoRecord
		}
		return Character{}, err
	}

	character.ID = id
	return character, nil
}

func (c *Character) deriveAttributes() map[string]int {
	tp := (c.Attributes["ko"] + c.Attributes["gr"]) / 10
	sta := c.Attributes["ma"]
	mp := c.Attributes["ma"] / 5

	res, _, err := dice.Roll("3d6kh3")
	if err != nil {
		return nil
	}
	luck := res.Int() * 5

	return map[string]int{
		"tp":   tp,
		"sta":  sta,
		"mp":   mp,
		"luck": luck,
	}
}
