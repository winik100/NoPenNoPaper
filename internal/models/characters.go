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
	Skills     CharacterSkills
	Items      []Item
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

type CharacterSkills struct {
	Anthropology      int
	Archaeology       int
	Driving           int
	LibraryResearch   int
	Accounting        int
	Charme            int
	CthulhuMythos     int
	Intimidate        int
	ElectricRepairs   int
	FirstAid          int
	Financials        int
	History           int
	Listening         int
	Concealing        int
	Climbing          int
	MechanicalRepairs int
	Medicine          int
	NaturalHistory    int
	Occultism         int
	Orientation       int
	PsychoAnalysis    int
	Psychology        int
	Law               int
	Horseriding       int
	Locks             int
	HeavyMachinery    int
	Swimming          int
	Jumping           int
	Tracking          int
	Persuasion        int
	Convincing        int
	Stealth           int
	DetectingSecrets  int
	Disguising        int
	Throwing          int
	Valuation         int
}

func DefaultCharacterSkills() CharacterSkills {
	var skills CharacterSkills
	skills.Anthropology = 1
	skills.Archaeology = 1
	skills.Driving = 20
	skills.LibraryResearch = 20
	skills.Accounting = 5
	skills.Charme = 15
	skills.CthulhuMythos = 0
	skills.Intimidate = 15
	skills.ElectricRepairs = 10
	skills.FirstAid = 30
	skills.Financials = 0
	skills.History = 5
	skills.Listening = 20
	skills.Concealing = 10
	skills.Climbing = 20
	skills.MechanicalRepairs = 10
	skills.Medicine = 1
	skills.NaturalHistory = 10
	skills.Occultism = 5
	skills.Orientation = 10
	skills.PsychoAnalysis = 1
	skills.Psychology = 10
	skills.Law = 5
	skills.Horseriding = 5
	skills.Locks = 1
	skills.HeavyMachinery = 1
	skills.Swimming = 20
	skills.Jumping = 20
	skills.Tracking = 10
	skills.Persuasion = 5
	skills.Convincing = 10
	skills.Stealth = 20
	skills.DetectingSecrets = 25
	skills.Disguising = 5
	skills.Throwing = 20
	skills.Valuation = 5
	return skills
}

func (s CharacterSkills) AsMap() map[string]int {
	return map[string]int{
		"Anthropologie":           s.Anthropology,
		"Archäologie":             s.Archaeology,
		"Autofahren":              s.Driving,
		"Bibliotheksnutzung":      s.LibraryResearch,
		"Buchführung":             s.Accounting,
		"Charme":                  s.Charme,
		"Cthulhu-Mythos":          s.CthulhuMythos,
		"Einschüchtern":           s.Intimidate,
		"Elektrische Reparaturen": s.ElectricRepairs,
		"Erste Hilfe":             s.FirstAid,
		"Finanzkraft":             s.Financials,
		"Geschichte":              s.History,
		"Horchen":                 s.Listening,
		"Kaschieren":              s.Concealing,
		"Klettern":                s.Climbing,
		"Mechanische Reparaturen": s.MechanicalRepairs,
		"Medizin":                 s.Medicine,
		"Naturkunde":              s.NaturalHistory,
		"Okkultismus":             s.Occultism,
		"Orientierung":            s.Orientation,
		"Psychoanalyse":           s.PsychoAnalysis,
		"Psychologie":             s.Psychology,
		"Rechtswesen":             s.Law,
		"Reiten":                  s.Horseriding,
		"Schließtechnik":          s.Locks,
		"Schweres Gerät":          s.HeavyMachinery,
		"Schwimmen":               s.Swimming,
		"Springen":                s.Jumping,
		"Spurensuche":             s.Tracking,
		"Überreden":               s.Persuasion,
		"Überzeugen":              s.Convincing,
		"Verborgen bleiben":       s.Stealth,
		"Verborgenes erkennen":    s.DetectingSecrets,
		"Verkleiden":              s.Disguising,
		"Werfen":                  s.Throwing,
		"Werte schätzen":          s.Valuation,
	}
}

type Item struct {
	ItemID      int
	Name        string
	Description string
	Count       int
}

type CharacterModelInterface interface {
	Insert(character Character, created_by int) (int, error)
	Get(characterId int) (Character, error)
	GetAllFrom(userId int) ([]Character, error)
	GetAll() ([]Character, error)
	AddItem(characterId int, item Item) error
	IncrementStat(character Character, stat string) (Character, error)
	DecrementStat(character Character, stat string) (Character, error)
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

	stmt = `INSERT INTO character_skills (character_id, anthropology, archaeology, driving, libraryResearch, accounting, charme, cthulhuMythos, intimidate, electricRepairs,
				firstAid, financials, history, listening, concealing, climbing, mechanicalRepairs, medicine, naturalHistory, occultism, orientation, psychoAnalysis, psychology, 
				law, horseriding, locks, heavyMachinery, swimming, jumping, tracking, persuasion, convincing, stealth, detectingSecrets, disguising, throwing, valuation)
			VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`
	_, err = tx.Exec(stmt, id, character.Skills.Anthropology, character.Skills.Archaeology, character.Skills.Driving, character.Skills.LibraryResearch, character.Skills.Accounting,
		character.Skills.Charme, character.Skills.CthulhuMythos, character.Skills.Intimidate, character.Skills.ElectricRepairs, character.Skills.FirstAid, character.Skills.Financials,
		character.Skills.History, character.Skills.Listening, character.Skills.Concealing, character.Skills.Climbing, character.Skills.MechanicalRepairs, character.Skills.Medicine,
		character.Skills.NaturalHistory, character.Skills.Occultism, character.Skills.Orientation, character.Skills.PsychoAnalysis, character.Skills.Psychology, character.Skills.Law,
		character.Skills.Horseriding, character.Skills.Locks, character.Skills.HeavyMachinery, character.Skills.Swimming, character.Skills.Jumping, character.Skills.Tracking,
		character.Skills.Persuasion, character.Skills.Convincing, character.Skills.Stealth, character.Skills.DetectingSecrets, character.Skills.Disguising, character.Skills.Throwing,
		character.Skills.Valuation)

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
	var sk CharacterSkills
	var items []Item

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

	stmt = `SELECT anthropology, archaeology, driving, libraryResearch, accounting, charme, cthulhuMythos, intimidate, electricRepairs,
				firstAid, financials, history, listening, concealing, climbing, mechanicalRepairs, medicine, naturalHistory, occultism, orientation, psychoAnalysis, psychology, 
				law, horseriding, locks, heavyMachinery, swimming, jumping, tracking, persuasion, convincing, stealth, detectingSecrets, disguising, throwing, valuation
			FROM character_skills
			WHERE character_id=?;`
	result = c.DB.QueryRow(stmt, characterId)
	err = result.Scan(&sk.Anthropology, &sk.Archaeology, &sk.Driving, &sk.LibraryResearch, &sk.Accounting, &sk.Charme, &sk.CthulhuMythos, &sk.Intimidate, &sk.ElectricRepairs, &sk.FirstAid,
		&sk.Financials, &sk.History, &sk.Listening, &sk.Concealing, &sk.Climbing, &sk.MechanicalRepairs, &sk.Medicine, &sk.NaturalHistory, &sk.Occultism, &sk.Orientation, &sk.PsychoAnalysis,
		&sk.Psychology, &sk.Law, &sk.Horseriding, &sk.Locks, &sk.HeavyMachinery, &sk.Swimming, &sk.Jumping, &sk.Tracking, &sk.Persuasion, &sk.Convincing, &sk.Stealth, &sk.DetectingSecrets,
		&sk.Disguising, &sk.Throwing, &sk.Valuation)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Character{}, ErrNoRecord
		}
		return Character{}, err
	}

	stmt = "SELECT name, description, cnt FROM items WHERE character_id=?;"
	rows, err := c.DB.Query(stmt, characterId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Character{}, ErrNoRecord
		}
		return Character{}, err
	}
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.Name, &item.Description, &item.Count)
		if err != nil {
			return Character{}, err
		}
		items = append(items, item)
	}

	return Character{ID: characterId, Info: info, Attributes: attr, Stats: stats, Skills: sk, Items: items}, nil
}

func (c *CharacterModel) GetAllFrom(userId int) ([]Character, error) {
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

func (c *CharacterModel) GetAll() ([]Character, error) {
	stmt := "SELECT id FROM characters;"
	rows, err := c.DB.Query(stmt)
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

func (c *CharacterModel) AddItem(characterId int, item Item) error {
	stmt := "INSERT INTO items (character_id, name, description, cnt) VALUES (?,?,?,?);"
	_, err := c.DB.Exec(stmt, characterId, item.Name, item.Description, item.Count)
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
