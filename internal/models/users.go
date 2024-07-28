package models

import (
	"database/sql"
	"errors"
	"time"
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
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) Insert(username, email, password string, isAdmin bool) (int, error) {
	admin, err := m.GetAdmin()
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return 0, err
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

func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
