package mocks

import (
	"github.com/winik100/NoPenNoPaper/internal/models"
)

// var MockUser = models.User{
// 	ID:             1,
// 	Name:           "Testnutzer",
// 	HashedPassword: "$2a$12$uK5Qivao7pieZMOZWtRTGubxPV3PgBf6ljFr3ACYtGPYZOrinx3ie", //"Klartext ole"
// }

type UserModel struct{}

func (m *UserModel) Insert(name, password string) error {
	return nil
}

func (m *UserModel) Authenticate(name, password string) (int, error) {
	if name == "Testnutzer" && password == "Klartext ole" {
		return 1, nil
	}
	return 0, models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id int) (bool, error) {
	if id == 1 || id == 2 {
		return true, nil
	}
	return false, nil
}

func (m *UserModel) GetRole(id int) (string, error) {
	if id == 1 {
		return "player", nil
	}
	if id == 2 {
		return "gm", nil
	}
	return "anonymous", nil
}
