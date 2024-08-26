package mocks

import (
	"github.com/winik100/NoPenNoPaper/internal/core"
	"github.com/winik100/NoPenNoPaper/internal/models"
)

var MockCharacterOtto = core.Character{
	ID:           1,
	Info:         mockInfo,
	Attributes:   mockAttributes,
	Skills:       mockSkills,
	CustomSkills: mockCustomSkills,
	Items:        mockItems,
	Notes:        mockNotes,
}

var MockCharacterViserys = core.Character{
	ID:   2,
	Info: mockInfo2,
}

var mockInfo = core.CharacterInfo{
	Name:       "Otto Hightower",
	Profession: "Lord von Oldtown",
	Age:        "65",
	Gender:     "männlich",
	Residence:  "Oldtown",
	Birthplace: "Oldtown",
}

var mockInfo2 = core.CharacterInfo{
	Name:       "Viserys Targaryen",
	Profession: "König von Westeros",
	Age:        "45",
	Gender:     "männlich",
	Residence:  "King's Landing",
	Birthplace: "King's Landing",
}

var mockAttributes = core.CharacterAttributes{
	ST: 40,
	GE: 50,
	MA: 50,
	KO: 50,
	ER: 70,
	BI: 60,
	GR: 60,
	IN: 80,
	BW: 6,
}

var mockSkills = core.Skills{
	Name:  []string{"Politik", "Intrige", "Manipulation"},
	Value: []int{70, 60, 60},
}

var mockCustomSkills = core.CustomSkills{
	Name:     []string{"Westerosi"},
	Value:    []int{50},
	Category: []string{"Muttersprache"},
}

var mockItems = core.Items{ItemId: []int{1}, Name: []string{"Hand-Brosche"}, Description: []string{"Brosche der Hand des Königs"}, Count: []int{1}}

var mockNotes = core.Notes{
	ID:   []int{1, 2},
	Text: []string{"Aegon ist blöde.", "Viserys war viel besser."}}

type CharacterModel struct{}

func (m *CharacterModel) Insert(character core.Character, created_by int) (int, error) {
	return 1, nil
}

func (m *CharacterModel) Get(characterId int) (core.Character, error) {
	if characterId == 1 {
		return MockCharacterOtto, nil
	}
	if characterId == 2 {
		return MockCharacterViserys, nil
	}
	return core.Character{}, models.ErrNoRecord
}

func (m *CharacterModel) Delete(characterId int) error {
	if characterId == 1 {
		MockCharacterOtto = core.Character{}
	}
	if characterId == 2 {
		MockCharacterViserys = core.Character{}
	}
	return nil
}

func (m *CharacterModel) GetAllFrom(userId int) ([]core.Character, error) {
	if userId == 1 {
		return []core.Character{MockCharacterOtto}, nil
	}
	if userId == 2 {
		return []core.Character{MockCharacterOtto, MockCharacterViserys}, nil
	}
	return nil, models.ErrNoRecord
}

func (m *CharacterModel) GetAll() ([]core.Character, error) {
	return []core.Character{MockCharacterOtto, MockCharacterViserys}, nil
}

func (m *CharacterModel) GetAvailableSkills() (core.Skills, error) {
	skills := core.Skills{Name: []string{"Politik", "Intrige", "Manipulation", "Schwertkampf", "Singen", "Tanzen"},
		Value: []int{10, 5, 5, 10, 20, 20}}
	return skills, nil
}

func (m *CharacterModel) AddSkill(characterId int, skill string, value int) error {
	return nil
}

func (m *CharacterModel) EditSkill(characterId int, skill string, newValue int) error {
	return nil
}

func (m *CharacterModel) AddCustomSkill(characterId int, Customkill string, category string, value int) error {
	return nil
}

func (m *CharacterModel) EditCustomSkill(characterId int, skill string, newValue int) error {
	return nil
}

func (m *CharacterModel) AddItem(characterId int, name, description string, count int) error {
	return nil
}

func (m *CharacterModel) EditItemCount(characterId, itemId, NewCount int) error {
	return nil
}

func (m *CharacterModel) DeleteItem(itemId int) error {
	if itemId == 1 {
		MockCharacterOtto.Items = core.Items{}
	}
	return nil
}

func (m *CharacterModel) AddNote(characterId int, text string) (int, error) {
	if characterId == 1 {
		return 2, nil
	}
	return 0, nil
}

func (m *CharacterModel) DeleteNote(noteId int) error {
	if noteId == 1 {
		MockCharacterOtto.Notes = core.Notes{ID: []int{2}, Text: []string{"Viserys war viel besser."}}
	}
	if noteId == 2 {
		MockCharacterOtto.Notes = core.Notes{ID: []int{1}, Text: []string{"Aegon ist blöde."}}
	}
	return nil
}

func (m *CharacterModel) IncrementStat(characterId int, stat string) (int, error) {
	character := MockCharacterOtto

	var updated int
	switch stat {
	case "TP":
		character.Stats.TP++
		updated = character.Stats.TP
	case "STA":
		character.Stats.STA++
		updated = character.Stats.STA
	case "MP":
		character.Stats.MP++
		updated = character.Stats.MP
	case "LUCK":
		character.Stats.LUCK++
		updated = character.Stats.LUCK
	}
	return updated, nil
}

func (m *CharacterModel) DecrementStat(characterId int, stat string) (int, error) {
	character := MockCharacterOtto

	var updated int
	switch stat {
	case "TP":
		character.Stats.TP--
		updated = character.Stats.TP
	case "STA":
		character.Stats.STA--
		updated = character.Stats.STA
	case "MP":
		character.Stats.MP--
		updated = character.Stats.MP
	case "LUCK":
		character.Stats.LUCK--
		updated = character.Stats.LUCK
	}
	return updated, nil
}
