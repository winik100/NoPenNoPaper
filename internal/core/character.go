package core

import (
	"slices"

	"github.com/justinian/dice"
)

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

func (character Character) AddableSkills(availableSkills Skills) Skills {
	var addableSkills Skills
	for i, sk := range availableSkills.Name {
		if !slices.Contains(character.Skills.Name, sk) {
			addableSkills.Name = append(addableSkills.Name, sk)
			addableSkills.Value = append(addableSkills.Value, availableSkills.Value[i])
		}
	}
	return addableSkills
}

func (character Character) DeriveStats() CharacterStats {
	tp := (character.Attributes.KO + character.Attributes.GR) / 10
	sta := character.Attributes.MA
	mp := character.Attributes.MA / 5

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
