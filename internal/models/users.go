package models

import (
	"database/sql"
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
