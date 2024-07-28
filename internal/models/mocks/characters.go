package mocks

import "github.com/winik100/NoPenNoPaper/internal/models"

var MockCharacter = models.Character{
	ID:           1,
	Info:         mockInfo,
	Attributes:   mockAttributes,
	Skills:       mockSkills,
	CustomSkills: mockCustomSkills,
	Items:        mockItems,
	Notes:        mockNotes,
}

var mockInfo = models.CharacterInfo{
	Name:       "Otto Hightower",
	Profession: "Lord von Oldtown",
	Age:        "65",
	Gender:     "männlich",
	Residence:  "Oldtown",
	Birthplace: "Oldtown",
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
		return MockCharacter, nil
	}
	return models.Character{}, models.ErrNoRecord
}

func (m *CharacterModel) GetAllFrom(userId int) ([]models.Character, error) {
	if userId == 1 {
		return []models.Character{MockCharacter}, nil
	}
	return nil, models.ErrNoRecord
}

func (m *CharacterModel) GetAll() ([]models.Character, error) {
	return []models.Character{MockCharacter}, nil
}

func (m *CharacterModel) GetAvailableSkills() (models.Skills, error) {
	skills := models.Skills{Name: []string{"Politik", "Intrige", "Manipulation", "Schwertkampf", "Singen", "Tanzen"},
		Value: []int{10, 5, 5, 10, 20, 20}}
	return skills, nil
}

func (m *CharacterModel) AddItem(characterId int, name, description string, count int) error {
	return nil
}

func (m *CharacterModel) DeleteItem(itemId int) error {
	return nil
}

func (m *CharacterModel) AddNote(characterId int, text string) error {
	return nil
}

func (m *CharacterModel) DeleteNote(noteId int) error {
	return nil
}

func (m *CharacterModel) IncrementStat(character models.Character, stat string) (models.Character, error) {
	var updated = character
	switch stat {
	case "TP":
		updated.Stats.TP++
	case "STA":
		updated.Stats.STA++
	case "MP":
		updated.Stats.MP++
	case "LUCK":
		updated.Stats.LUCK++
	}
	return updated, nil
}

func (m *CharacterModel) DecrementStat(character models.Character, stat string) (models.Character, error) {
	var updated = character
	switch stat {
	case "TP":
		updated.Stats.TP--
	case "STA":
		updated.Stats.STA--
	case "MP":
		updated.Stats.MP--
	case "LUCK":
		updated.Stats.LUCK--
	}
	return updated, nil
}
