package models

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	HashedPassword string
}

type UserModelInterface interface {
	Insert(name, password string) error
	Exists(id int) (bool, error)
}

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) Insert(name, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := "INSERT INTO users (name, hashed_password) VALUES (?,?);"
	_, err = u.DB.Exec(stmt, name, hashedPassword)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserModel) Exists(id int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?);"
	err := u.DB.QueryRow(stmt, id).Scan(&exists)

	return exists, err
}
