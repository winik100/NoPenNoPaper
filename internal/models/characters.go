package models

import (
	"database/sql"
	"errors"

	"github.com/justinian/dice"
)

type Character struct {
	ID         int
	Info       CharacterInfo
	Attributes CharacterAttributes
	Stats      CharacterStats
}

type CharacterInfo struct {
	Name       string
	Profession string
	Age        string
	Gender     string
	Residence  string
	Birthplace string
}

func (ci *CharacterInfo) AsMap() map[string]string {
	return map[string]string{
		"name":       ci.Name,
		"profession": ci.Profession,
		"age":        ci.Age,
		"gender":     ci.Gender,
		"residence":  ci.Residence,
		"birthplace": ci.Birthplace,
	}
}

type CharacterAttributes struct {
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

func (a *CharacterAttributes) AsMap() map[string]int {
	return map[string]int{
		"st": a.ST,
		"ge": a.GE,
		"ma": a.MA,
		"ko": a.KO,
		"er": a.ER,
		"BI": a.BI,
		"GR": a.GR,
		"IN": a.IN,
		"BW": a.BW,
	}
}

type CharacterStats struct {
	TP   int
	STA  int
	MP   int
	LUCK int
}

func (st *CharacterStats) asMap() map[string]int {
	return map[string]int{
		"tp":   st.TP,
		"sta":  st.STA,
		"mp":   st.MP,
		"luck": st.LUCK,
	}
}

type CharacterModelInterface interface {
	Insert(character Character, created_by int) (int, error)
	Get(characterId int) (Character, error)
	GetAll(userId int) ([]Character, error)
}

type CharacterModel struct {
	DB *sql.DB
}

func (c *Character) deriveStats() CharacterStats {
	tp := (c.Attributes.KO + c.Attributes.GR) / 10
	sta := c.Attributes.MA
	mp := c.Attributes.MA / 5

	res, _, err := dice.Roll("3d6kh3")
	if err != nil {
		return CharacterStats{}
	}
	luck := res.Int() * 5

	return CharacterStats{
		TP:   tp,
		STA:  sta,
		MP:   mp,
		LUCK: luck,
	}
}

func (c *CharacterModel) Insert(character Character, created_by int) (int, error) {
	tx, err := c.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	stmt := "INSERT INTO characters (created_by) VALUES (?);"
	result, err := tx.Exec(stmt, created_by)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	stmt = "INSERT INTO character_info (character_id, name, profession, age, gender, residence, birthplace) VALUES (?,?,?,?,?,?,?);"
	_, err = tx.Exec(stmt, id, character.Info.Name, character.Info.Profession, character.Info.Age, character.Info.Gender, character.Info.Residence, character.Info.Birthplace)
	if err != nil {
		return 0, err
	}

	stmt = "INSERT INTO character_attributes (character_id, st, ge, ma, ko, er, bi, gr, i, bw) VALUES (?,?,?,?,?,?,?,?,?,?);"
	_, err = tx.Exec(stmt, id, character.Attributes.ST, character.Attributes.GE, character.Attributes.MA,
		character.Attributes.KO, character.Attributes.ER, character.Attributes.BI,
		character.Attributes.GR, character.Attributes.IN, character.Attributes.BW)
	if err != nil {
		return 0, err
	}

	stats := character.deriveStats()
	stmt = "INSERT INTO character_stats (character_id, tp, sta, mp, luck) VALUES (?,?,?,?,?);"
	_, err = tx.Exec(stmt, id, stats.TP, stats.STA, stats.MP, stats.LUCK)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (c *CharacterModel) Get(characterId int) (Character, error) {
	var info CharacterInfo
	var attr CharacterAttributes
	var stats CharacterStats

	stmt := "SELECT name, profession, age, gender, residence, birthplace FROM character_info WHERE character_id=?;"
	result := c.DB.QueryRow(stmt, characterId)
	err := result.Scan(&info.Name, &info.Profession, &info.Age, &info.Gender, &info.Residence, &info.Birthplace)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Character{}, ErrNoRecord
		}
		return Character{}, err
	}

	stmt = "SELECT st, ge, ma, ko, er, bi, gr, i, bw FROM character_attributes WHERE character_id=?;"
	result = c.DB.QueryRow(stmt, characterId)
	err = result.Scan(&attr.ST, &attr.GE, &attr.MA, &attr.KO, &attr.ER, &attr.BI, &attr.GR, &attr.IN, &attr.BW)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Character{}, ErrNoRecord
		}
		return Character{}, err
	}

	stmt = "SELECT tp, sta, mp, luck FROM character_stats WHERE character_id=?;"
	result = c.DB.QueryRow(stmt, characterId)
	err = result.Scan(&stats.TP, &stats.STA, &stats.MP, &stats.LUCK)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Character{}, ErrNoRecord
		}
		return Character{}, err
	}

	return Character{ID: characterId, Info: info, Attributes: attr, Stats: stats}, nil
}

func (c *CharacterModel) GetAll(userId int) ([]Character, error) {
	stmt := "SELECT id FROM characters WHERE created_by=?;"
	rows, err := c.DB.Query(stmt, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var characterIds []int
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		characterIds = append(characterIds, id)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	var characters []Character
	for _, id := range characterIds {
		character, err := c.Get(id)
		if err != nil {
			return nil, err
		}
		characters = append(characters, character)
	}

	return characters, nil
}
