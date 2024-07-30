package models

import (
	"database/sql"
	"errors"
	"slices"

	"github.com/justinian/dice"
)

type CharacterModelInterface interface {
	Insert(character Character, created_by int) (int, error)
	Get(characterId int) (Character, error)
	GetAllFrom(userId int) ([]Character, error)
	GetAll() ([]Character, error)
	Delete(characterId int) error
	GetAvailableSkills() (Skills, error)
	AddSkill(characterId int, skill string, value int) error
	EditSkill(characterId int, skill string, newValue int) error
	EditCustomSkill(characterId int, skill string, newValue int) error
	AddItem(characterId int, name, description string, count int) error
	DeleteItem(itemId int) error
	AddNote(characterId int, text string) (int, error)
	DeleteNote(noteId int) error
	IncrementStat(character Character, stat string) (Character, error)
	DecrementStat(character Character, stat string) (Character, error)
}

type CharacterModel struct {
	DB *sql.DB
}

type Character struct {
	ID           int
	Info         CharacterInfo
	Attributes   CharacterAttributes
	Stats        CharacterStats
	Skills       Skills
	CustomSkills CustomSkills
	Items        Items
	Notes        Notes
}

type CharacterInfo struct {
	Name       string
	Profession string
	Age        string
	Gender     string
	Residence  string
	Birthplace string
}

func (ci CharacterInfo) AsMap() map[string]string {
	return map[string]string{
		"Name":       ci.Name,
		"Beruf":      ci.Profession,
		"Alter":      ci.Age,
		"Geschlecht": ci.Gender,
		"Wohnort":    ci.Residence,
		"Geburtsort": ci.Birthplace,
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

func (a CharacterAttributes) AsMap() map[string]int {
	return map[string]int{
		"ST": a.ST,
		"GE": a.GE,
		"MA": a.MA,
		"KO": a.KO,
		"ER": a.ER,
		"BI": a.BI,
		"GR": a.GR,
		"IN": a.IN,
		"BW": a.BW,
	}
}

func (a CharacterAttributes) OrderedKeys() []string {
	return []string{"ST", "GE", "MA", "KO", "ER", "BI", "GR", "IN", "BW"}
}

type CharacterStats struct {
	MaxTP   int
	TP      int
	MaxSTA  int
	STA     int
	MaxMP   int
	MP      int
	MaxLUCK int
	LUCK    int
}

func (st CharacterStats) OrderedKeysCurrent() []string {
	return []string{"TP", "STA", "MP", "LUCK"}
}

func (st CharacterStats) CurrentAsMap() map[string]int {
	return map[string]int{
		"TP":   st.TP,
		"STA":  st.STA,
		"MP":   st.MP,
		"LUCK": st.LUCK,
	}
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
		MaxTP:   tp,
		TP:      tp,
		MaxSTA:  sta,
		STA:     sta,
		MaxMP:   mp,
		MP:      mp,
		MaxLUCK: luck,
		LUCK:    luck,
	}
}

// workaround due to gorillas.schema not being able to parse into slices of structs
type Skills struct {
	Name  []string
	Value []int
}

type CustomSkills struct {
	Category []string
	Name     []string
	Value    []int
}

func (character Character) AllSkills() Skills {
	allNames := append(character.Skills.Name, character.CustomSkills.Name...)
	allValues := append(character.Skills.Value, character.CustomSkills.Value...)

	slices.Sort(allNames)
	slices.Sort(allValues)
	return Skills{Name: allNames, Value: allValues}
}

type Items struct {
	ItemId      []int
	Name        []string
	Description []string
	Count       []int
}

type Notes struct {
	ID   []int
	Text []string
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
	stmt = "INSERT INTO character_stats (character_id, maxtp, tp, maxsta, sta, maxmp, mp, maxluck, luck) VALUES (?,?,?,?,?,?,?,?,?);"
	_, err = tx.Exec(stmt, id, stats.TP, stats.TP, stats.STA, stats.STA, stats.MP, stats.MP, stats.LUCK, stats.LUCK)
	if err != nil {
		return 0, err
	}

	for i, customSkill := range character.CustomSkills.Name {
		var exists bool
		stmt = "SELECT EXISTS(SELECT true FROM custom_skills WHERE name=? AND category=?);"
		err = tx.QueryRow(stmt, customSkill, character.CustomSkills.Category[i]).Scan(&exists)
		if err != nil {
			return 0, err
		}
		if !exists {
			stmt = "INSERT INTO custom_skills (name, category, default_value) VALUES (?,?,?);"
			_, err = tx.Exec(stmt, customSkill, character.CustomSkills.Category[i], DefaultForCategory(character.CustomSkills.Category[i]))
			if err != nil {
				return 0, err
			}
		}

		stmt = "INSERT INTO character_custom_skills (character_id, custom_skill_name, value) VALUES (?,?,?);"
		_, err = tx.Exec(stmt, id, customSkill, character.CustomSkills.Value[i])
		if err != nil {
			return 0, err
		}
	}

	for i, skill := range character.Skills.Name {
		stmt = "INSERT INTO character_skills (character_id, skill_name, value) VALUES (?,?,?);"
		_, err = tx.Exec(stmt, id, skill, character.Skills.Value[i])
		if err != nil {
			return 0, err
		}
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
	var skills Skills
	var customSkills CustomSkills
	var items Items
	var notes Notes

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

	stmt = "SELECT maxtp, tp, maxsta, sta, maxmp, mp, maxluck, luck FROM character_stats WHERE character_id=?;"
	result = c.DB.QueryRow(stmt, characterId)
	err = result.Scan(&stats.MaxTP, &stats.TP, &stats.MaxSTA, &stats.STA, &stats.MaxMP, &stats.MP, &stats.MaxLUCK, &stats.LUCK)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Character{}, ErrNoRecord
		}
		return Character{}, err
	}

	stmt = "SELECT skill_name, value FROM character_skills WHERE character_id=?;"
	rows, err := c.DB.Query(stmt, characterId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Character{}, ErrNoRecord
		}
		return Character{}, err
	}

	var skillsName []string
	var skillsValue []int
	for rows.Next() {
		var name string
		var value int

		err = rows.Scan(&name, &value)
		if err != nil {
			return Character{}, err
		}
		skillsName = append(skillsName, name)
		skillsValue = append(skillsValue, value)
	}
	skills.Name = skillsName
	skills.Value = skillsValue

	stmt = "SELECT custom_skill_name, value FROM character_custom_skills WHERE character_id=?;"
	rows, err = c.DB.Query(stmt, characterId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Character{}, ErrNoRecord
		}
		return Character{}, err
	}
	var customSkillsName []string
	var customSkillsValue []int
	for rows.Next() {
		var name string
		var value int

		err = rows.Scan(&name, &value)
		if err != nil {
			return Character{}, err
		}
		customSkillsName = append(customSkillsName, name)
		customSkillsValue = append(customSkillsValue, value)
	}
	customSkills.Name = customSkillsName
	customSkills.Value = customSkillsValue

	stmt = "SELECT item_id, name, description, cnt FROM items WHERE character_id=?;"
	rows, err = c.DB.Query(stmt, characterId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Character{}, ErrNoRecord
		}
		return Character{}, err
	}
	for rows.Next() {
		var id, count int
		var name, description string
		err = rows.Scan(&id, &name, &description, &count)
		if err != nil {
			return Character{}, err
		}
		items.ItemId = append(items.ItemId, id)
		items.Name = append(items.Name, name)
		items.Description = append(items.Description, description)
		items.Count = append(items.Count, count)
	}

	stmt = "SELECT note_id, text FROM notes WHERE character_id=?;"
	rows, err = c.DB.Query(stmt, characterId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Character{}, ErrNoRecord
		}
		return Character{}, err
	}
	for rows.Next() {
		var id int
		var text string
		err = rows.Scan(&id, &text)
		if err != nil {
			return Character{}, err
		}
		notes.ID = append(notes.ID, id)
		notes.Text = append(notes.Text, text)
	}

	return Character{ID: characterId, Info: info, Attributes: attr, Stats: stats, Skills: skills, CustomSkills: customSkills, Items: items, Notes: notes}, nil
}

func (c *CharacterModel) Delete(characterId int) error {
	stmt := "DELETE FROM characters WHERE id=?;"
	_, err := c.DB.Exec(stmt, characterId)
	if err != nil {
		return err
	}
	return nil
}

func (c *CharacterModel) GetAllFrom(userId int) ([]Character, error) {
	stmt := "SELECT id FROM characters WHERE created_by=?;"
	rows, err := c.DB.Query(stmt, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
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

func (c *CharacterModel) GetAll() ([]Character, error) {
	stmt := "SELECT id FROM characters;"
	rows, err := c.DB.Query(stmt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
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

func (c *CharacterModel) GetAvailableSkills() (Skills, error) {
	var skillsName []string
	var skillsValue []int
	stmt := "SELECT name, default_value FROM skills;"
	rows, err := c.DB.Query(stmt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Skills{}, ErrNoRecord
		}
		return Skills{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var value int
		err = rows.Scan(&name, &value)
		if err != nil {
			return Skills{}, err
		}
		skillsName = append(skillsName, name)
		skillsValue = append(skillsValue, value)
	}
	return Skills{Name: skillsName, Value: skillsValue}, nil
}

func (c *CharacterModel) AddSkill(characterId int, skill string, value int) error {
	stmt := "INSERT INTO character_skills (character_id, skill_name, value) VALUES (?,?,?);"
	_, err := c.DB.Exec(stmt, characterId, skill, value)
	if err != nil {
		return err
	}
	return nil
}

func (c *CharacterModel) EditSkill(characterId int, skill string, newValue int) error {
	stmt := "UPDATE character_skills SET value=? WHERE character_id=? AND skill_name=?;"
	_, err := c.DB.Exec(stmt, newValue, characterId, skill)
	if err != nil {
		return err
	}
	return nil
}

func (c *CharacterModel) EditCustomSkill(characterId int, skill string, newValue int) error {
	stmt := "UPDATE character_custom_skills SET value=? WHERE character_id=? AND custom_skill_name=?;"
	_, err := c.DB.Exec(stmt, newValue, characterId, skill)
	if err != nil {
		return err
	}
	return nil
}

func (c *CharacterModel) AddItem(characterId int, name, description string, count int) error {
	stmt := "INSERT INTO items (character_id, name, description, cnt) VALUES (?,?,?,?);"
	_, err := c.DB.Exec(stmt, characterId, name, description, count)
	if err != nil {
		return err
	}
	return nil
}

func (c *CharacterModel) DeleteItem(itemId int) error {
	stmt := "DELETE FROM items WHERE item_id=?;"
	_, err := c.DB.Exec(stmt, itemId)
	if err != nil {
		return err
	}
	return nil
}

func (c *CharacterModel) AddNote(characterId int, text string) (int, error) {
	stmt := "INSERT INTO notes (character_id, text) VALUES (?,?);"
	res, err := c.DB.Exec(stmt, characterId, text)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (c *CharacterModel) DeleteNote(noteId int) error {
	stmt := "DELETE FROM notes WHERE note_id=?;"
	_, err := c.DB.Exec(stmt, noteId)
	if err != nil {
		return err
	}
	return nil
}

func (c *CharacterModel) IncrementStat(character Character, stat string) (Character, error) {
	var stmt string
	switch stat {
	case "TP":
		stmt = "UPDATE character_stats SET tp=? WHERE character_id=?;"
		_, err := c.DB.Exec(stmt, character.Stats.TP+1, character.ID)
		if err != nil {
			return character, err
		}
		character.Stats.TP = character.Stats.TP + 1
	case "STA":
		stmt = "UPDATE character_stats SET sta=? WHERE character_id=?;"
		_, err := c.DB.Exec(stmt, character.Stats.STA+1, character.ID)
		if err != nil {
			return character, err
		}
		character.Stats.STA = character.Stats.STA + 1
	case "MP":
		stmt = "UPDATE character_stats SET mp=? WHERE character_id=?;"
		_, err := c.DB.Exec(stmt, character.Stats.MP+1, character.ID)
		if err != nil {
			return character, err
		}
		character.Stats.MP = character.Stats.MP + 1
	case "LUCK":
		stmt = "UPDATE character_stats SET luck=? WHERE character_id=?;"
		_, err := c.DB.Exec(stmt, character.Stats.LUCK+1, character.ID)
		if err != nil {
			return character, err
		}
		character.Stats.LUCK = character.Stats.LUCK + 1
	}
	return character, nil
}

func (c *CharacterModel) DecrementStat(character Character, stat string) (Character, error) {
	var stmt string
	switch stat {
	case "TP":
		stmt = "UPDATE character_stats SET tp=? WHERE character_id=?;"
		_, err := c.DB.Exec(stmt, character.Stats.TP-1, character.ID)
		if err != nil {
			return character, err
		}
		character.Stats.TP = character.Stats.TP - 1
	case "STA":
		stmt = "UPDATE character_stats SET sta=? WHERE character_id=?;"
		_, err := c.DB.Exec(stmt, character.Stats.STA-1, character.ID)
		if err != nil {
			return character, err
		}
		character.Stats.STA = character.Stats.STA - 1
	case "MP":
		stmt = "UPDATE character_stats SET mp=? WHERE character_id=?;"
		_, err := c.DB.Exec(stmt, character.Stats.MP-1, character.ID)
		if err != nil {
			return character, err
		}
		character.Stats.MP = character.Stats.MP - 1
	case "LUCK":
		stmt = "UPDATE character_stats SET luck=? WHERE character_id=?;"
		_, err := c.DB.Exec(stmt, character.Stats.LUCK-1, character.ID)
		if err != nil {
			return character, err
		}
		character.Stats.LUCK = character.Stats.LUCK - 1
	}
	return character, nil
}

func DefaultForCategory(category string) int {
	switch category {
	case "Muttersprache":
		return 50
	case "Fremdsprache":
		return 1
	case "Handwerk":
		return 5
	case "Naturwissenschaft":
		return 1
	case "Steuern":
		return 1
	case "Überlebenskunst":
		return 10
	}
	return -1
}
