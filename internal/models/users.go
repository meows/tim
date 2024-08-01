package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	CreatedAt time.Time `db:"created_at"`
	LastLogin time.Time `db:"last_login"`
	ID        string    `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	Email     string    `db:"email"`
	IsAdmin   bool      `db:"is_admin"`
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

func (m *UserModel) Insert(username, email, password string, isAdmin bool) (string, error) {
	admin, err := m.GetAdmin()
	if err != nil {
		fmt.Println(err, "here is error")
		if !errors.Is(err, ErrNoRecord) {
			return "", ErrNoRecord
		}
	}

	if isAdmin && admin != nil {
		return "", ErrDuplicateAdmin
	}

	stmt := `INSERT INTO users (id, username, email, password, is_admin, last_login) VALUES(?, ?, ?, ?, ?, ?)`

	id := uuid.New().String()
	_, err = m.DB.Exec(stmt, id, username, strings.ToLower(email), password, isAdmin, time.Now().UTC())
	if err != nil {
		return "", err
	}

	//	insertId, err := result.LastInsertId()
	//	if err != nil {
	//		return "", err
	//	}

	return id, nil
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
