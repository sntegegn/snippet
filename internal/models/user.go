package models

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {

	hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created) 
			VALUE (?, ?, ?, UTC_TIMESTAMP())`
	_, err = m.DB.Exec(stmt, name, email, string(hashed_password))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashed_password []byte
	stmt := `SELECT id, hashed_password FROM users WHERE email=?`

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashed_password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword(hashed_password, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}
	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool

	stmt := `SELECT EXISTS(SELECT true FROM users WHERE id=?)`

	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}
