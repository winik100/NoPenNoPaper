package models

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/winik100/NoPenNoPaper/internal/core"
	"golang.org/x/crypto/bcrypt"
)

type UserModelInterface interface {
	Insert(name, password string) error
	Get(name string) (core.User, error)
	Authenticate(name, password string) (int, error)
	Exists(userId int) (bool, error)
	GetRole(id int) (string, error)
	AddMaterial(fileName string, uploadedBy int) error
}

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) Insert(name, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := "INSERT INTO users (name, hashed_password, role) VALUES (?,?,?);"
	_, err = u.DB.Exec(stmt, name, hashedPassword, "player")
	if err != nil {
		return err
	}
	return nil
}

func (u *UserModel) Get(name string) (core.User, error) {
	stmt := "SELECT id, hashed_password FROM users WHERE name=?;"
	row := u.DB.QueryRow(stmt, name)

	var user core.User
	err := row.Scan(&user.ID, &user.HashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.User{}, ErrNoRecord
		}
		return core.User{}, err
	}
	user.Name = name
	return user, nil
}

func (u *UserModel) Authenticate(name, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password FROM users where name=?;"
	err := u.DB.QueryRow(stmt, name).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (u *UserModel) Exists(userId int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id=?);"
	err := u.DB.QueryRow(stmt, userId).Scan(&exists)

	return exists, err
}

func (u *UserModel) GetRole(id int) (string, error) {
	stmt := "SELECT role FROM users WHERE id=?;"

	var role string
	err := u.DB.QueryRow(stmt, id).Scan(&role)
	if err != nil {
		return "", err
	}

	return role, nil
}

func (u *UserModel) AddMaterial(fileName string, uploadedBy int) error {
	stmt := "INSERT INTO materials (file_name, uploaded_by) VALUES (?,?);"

	_, err := u.DB.Exec(stmt, fileName, uploadedBy)
	var mysqlErr *mysql.MySQLError
	if err != nil {
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return ErrDuplicateFileName
		}
		return err
	}
	return nil
}
