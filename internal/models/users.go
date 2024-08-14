package models

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/winik100/NoPenNoPaper/internal/core"
	"golang.org/x/crypto/bcrypt"
)

type UserModelInterface interface {
	Insert(name, password string) (int, error)
	Get(name string) (core.User, error)
	Authenticate(name, password string) (int, error)
	Exists(userName string) (bool, error)
	AddMaterial(title string, fileName string, uploadedBy int) error
	DeleteMaterial(fileName string, uploadedBy int) error
}

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) Insert(name, password string) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err
	}

	stmt := "INSERT INTO users (name, hashed_password, role) VALUES (?,?,?);"
	//placeholder role
	res, err := u.DB.Exec(stmt, name, hashedPassword, "player")
	if err != nil {
		return 0, err
	}
	userId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(userId), nil
}

func (u *UserModel) Get(name string) (core.User, error) {
	stmt := "SELECT id, hashed_password, role FROM users WHERE name=?;"
	row := u.DB.QueryRow(stmt, name)

	var user core.User
	err := row.Scan(&user.ID, &user.HashedPassword, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.User{}, ErrNoRecord
		}
		return core.User{}, err
	}

	// stmt = `SELECT u.id, u.hashed_password, u.role, m.title, m.file_name FROM users AS u LEFT JOIN materials AS m ON u.id = m.uploaded_by WHERE u.name = ?;`

	stmt = "SELECT title, file_name FROM materials WHERE uploaded_by=?;"
	rows, err := u.DB.Query(stmt, user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.User{}, ErrNoRecord
		}
		return core.User{}, err
	}

	var titles, fileNames []string
	for rows.Next() {
		var title string
		var fileName string
		err = rows.Scan(&title, &fileName)
		if err != nil {
			return core.User{}, err
		}
		titles = append(titles, title)
		fileNames = append(fileNames, fileName)
	}

	user.Materials = core.Materials{
		Title:    titles,
		FileName: fileNames,
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

func (u *UserModel) Exists(userName string) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM users WHERE name=?);"
	err := u.DB.QueryRow(stmt, userName).Scan(&exists)

	return exists, err
}

func (u *UserModel) AddMaterial(title string, fileName string, uploadedBy int) error {
	stmt := "INSERT INTO materials (title, file_name, uploaded_by) VALUES (?, ?,?);"

	_, err := u.DB.Exec(stmt, title, fileName, uploadedBy)
	var mysqlErr *mysql.MySQLError
	if err != nil {
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return ErrDuplicateFileName
		}
		return err
	}
	return nil
}

func (u *UserModel) DeleteMaterial(fileName string, uploadedBy int) error {
	stmt := "DELETE FROM materials WHERE file_name = ? AND uploaded_by = ?"

	_, err := u.DB.Exec(stmt, fileName, uploadedBy)
	if err != nil {
		return err
	}
	return nil
}
