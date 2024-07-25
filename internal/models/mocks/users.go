package mocks

import "github.com/winik100/NoPenNoPaper/internal/models"

var mockUser = models.User{
	ID:             1,
	Name:           "Testnutzer",
	HashedPassword: "$2a$12$uK5Qivao7pieZMOZWtRTGubxPV3PgBf6ljFr3ACYtGPYZOrinx3ie", //"Klartext ole"
}

type UserModel struct{}

func (m *UserModel) Insert(name, password string) error {
	return nil
}

func (m *UserModel) Authenticate(name, password string) (int, error) {
	return 1, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return true, nil
}

func (m *UserModel) GetRole(id int) (string, error) {
	return "player", nil
}
