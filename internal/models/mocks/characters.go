package mocks

import (
	"github.com/winik100/NoPenNoPaper/internal/models"
)

var MockCharacterOtto = models.Character{
	ID:           1,
	Info:         mockInfo,
	Attributes:   mockAttributes,
	Skills:       mockSkills,
	CustomSkills: mockCustomSkills,
	Items:        mockItems,
	Notes:        mockNotes,
}

var MockCharacterViserys = models.Character{
	ID:   2,
	Info: mockInfo2,
}

var mockInfo = models.CharacterInfo{
	Name:       "Otto Hightower",
	Profession: "Lord von Oldtown",
	Age:        "65",
	Gender:     "männlich",
	Residence:  "Oldtown",
	Birthplace: "Oldtown",
}

var mockInfo2 = models.CharacterInfo{
	Name:       "Viserys Targaryen",
	Profession: "König von Westeros",
	Age:        "45",
	Gender:     "männlich",
	Residence:  "King's Landing",
	Birthplace: "King's Landing",
}

var mockAttributes = models.CharacterAttributes{
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

var mockSkills = models.Skills{
	Name:  []string{"Politik", "Intrige", "Manipulation"},
	Value: []int{70, 60, 60},
}

var mockCustomSkills = models.CustomSkills{
	Name:     []string{"Westerosi"},
	Value:    []int{50},
	Category: []string{"Muttersprache"},
}

var mockItems = models.Items{ItemId: []int{1}, Name: []string{"Hand-Brosche"}, Description: []string{"Brosche der Hand des Königs"}, Count: []int{1}}

var mockNotes = models.Notes{
	ID:   []int{1, 2},
	Text: []string{"Aegon ist blöde.", "Viserys war viel besser."}}

type CharacterModel struct{}

func (m *CharacterModel) Insert(character models.Character, created_by int) (int, error) {
	return 1, nil
}

func (m *CharacterModel) Get(characterId int) (models.Character, error) {
	if characterId == 1 {
		return MockCharacterOtto, nil
	}
	if characterId == 2 {
		return MockCharacterViserys, nil
	}
	return models.Character{}, models.ErrNoRecord
}

func (m *CharacterModel) Delete(characterId int) error {
	if characterId == 1 {
		MockCharacterOtto = models.Character{}
	}
	if characterId == 2 {
		MockCharacterViserys = models.Character{}
	}
	return nil
}

func (m *CharacterModel) GetAllFrom(userId int) ([]models.Character, error) {
	if userId == 1 {
		return []models.Character{MockCharacterOtto}, nil
	}
	return nil, models.ErrNoRecord
}

func (m *CharacterModel) GetAll() ([]models.Character, error) {
	return []models.Character{MockCharacterOtto, MockCharacterViserys}, nil
}

func (m *CharacterModel) GetAvailableSkills() (models.Skills, error) {
	skills := models.Skills{Name: []string{"Politik", "Intrige", "Manipulation", "Schwertkampf", "Singen", "Tanzen"},
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
		MockCharacterOtto.Items = models.Items{}
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
		MockCharacterOtto.Notes = models.Notes{ID: []int{2}, Text: []string{"Viserys war viel besser."}}
	}
	if noteId == 2 {
		MockCharacterOtto.Notes = models.Notes{ID: []int{1}, Text: []string{"Aegon ist blöde."}}
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
