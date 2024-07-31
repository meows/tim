package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"password" db:"password"`
	Email     string    `json:"email" db:"email"`
	IsAdmin   bool      `json:"is_admin" db:"is_admin"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	LastLogin time.Time `json:"last_login" db:"last_login"`
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) GetAdmin() (*User, error) {
	var user User
	// Get single user with is_admin = true
	query := `SELECT * FROM users WHERE is_admin = true`

	err := m.DB.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.IsAdmin, &user.CreatedAt, &user.LastLogin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) Insert(username, email, password string, isAdmin bool) (int, error) {
	admin, err := m.GetAdmin()
	if err != nil {
		fmt.Println(err, "here is error")
		if !errors.Is(err, ErrNoRecord) {
			return 0, ErrNoRecord
		}
	}

	if isAdmin && admin != nil {
		return 0, ErrDuplicateAdmin
	}

	stmt := `INSERT INTO users (username, email, password, is_admin, last_login) VALUES(?, ?, ?, ?, ?)`

	result, err := m.DB.Exec(stmt, username, email, password, isAdmin, time.Now().UTC())
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	var user User
	query := `SELECT * FROM users WHERE email = ?`
	err := m.DB.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.IsAdmin, &user.CreatedAt, &user.LastLogin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := `SELECT id, password FROM users WHERE email = ?`

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
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

	return 0, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
