package mocks

import (
	"github.com/winik100/NoPenNoPaper/internal/core"
	"github.com/winik100/NoPenNoPaper/internal/models"
)

var MockPlayer = core.User{
	ID:             1,
	Name:           "Testnutzer",
	HashedPassword: "$2a$12$uK5Qivao7pieZMOZWtRTGubxPV3PgBf6ljFr3ACYtGPYZOrinx3ie", //"Klartext ole"
	Role:           "player",
}

var MockGM = core.User{
	ID:             2,
	Name:           "Test-GM",
	HashedPassword: "$2a$12$uK5Qivao7pieZMOZWtRTGubxPV3PgBf6ljFr3ACYtGPYZOrinx3ie",
	Role:           "gm",
}

type UserModel struct{}

func (m *UserModel) Insert(name, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Get(name string) (core.User, error) {
	if name == MockPlayer.Name {
		return MockPlayer, nil
	}
	if name == MockGM.Name {
		return MockGM, nil
	}
	return core.User{}, models.ErrNoRecord
}

func (m *UserModel) Delete(name string) error {
	return nil
}

func (m *UserModel) Authenticate(name, password string) (int, error) {
	if name == MockPlayer.Name && password == "Klartext ole" {
		return 1, nil
	}
	if name == MockGM.Name && password == "Klartext ole" {
		return 2, nil
	}
	return 0, models.ErrInvalidCredentials
}

func (m *UserModel) Exists(userName string) (bool, error) {
	if userName == MockPlayer.Name || userName == MockGM.Name {
		return true, nil
	}
	return false, nil
}

func (m *UserModel) AddMaterial(title string, fileName string, uploadedBy int) error {
	return nil
}

func (m *UserModel) DeleteMaterial(fileName string, uploadedBy int) error {
	return nil
}
