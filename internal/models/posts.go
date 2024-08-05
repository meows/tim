package models

import (
	"database/sql"
	"errors"
	"time"
)

type Post struct {
	ID          int       `json:"id" db:"id"`
	AuthorID    int       `json:"author_id" db:"author_id"`
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content" db:"content"`
	ContentHTML string    `json:"content_html"`
	Private     bool      `json:"private" db:"private"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) Insert(title, content string, private bool, authorID int) (int, error) {
	stmt := `INSERT INTO posts (title, content, author_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	result, err := m.DB.Exec(stmt, title, content, authorID, time.Now(), time.Now())
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *PostModel) Get(id int) (Post, error) {
	var p Post
	stmt := `SELECT id, title, content, private, created_at, updated_at, author_id FROM posts WHERE id = ?`

	err := m.DB.QueryRow(stmt, id).Scan(&p.ID, &p.Title, &p.Content, &p.Private, &p.CreatedAt, &p.UpdatedAt, &p.AuthorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Post{}, ErrNoRecord
		} else {
			return Post{}, err
		}
	}

	return p, nil
}

func (m *PostModel) GetPostByTitle(title string) (Post, error) {
	var p Post
	stmt := `SELECT id, title, content, private, created_at, updated_at, author_id FROM posts WHERE title = ?`

	err := m.DB.QueryRow(stmt, title).Scan(&p.ID, &p.Title, &p.Content, &p.Private, &p.CreatedAt, &p.UpdatedAt, &p.AuthorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Post{}, ErrNoRecord
		} else {
			return Post{}, err
		}
	}

	return p, nil
}

func (m *PostModel) Latest(includePrivatePosts bool) ([]Post, error) {
	var stmt string
	if includePrivatePosts {
		stmt = `SELECT id, title, content, private, created_at, updated_at, author_id FROM posts ORDER BY created_at DESC LIMIT 10`
	} else {
		stmt = `SELECT id, title, content, private, created_at, updated_at, author_id FROM posts WHERE private=false ORDER BY created_at DESC LIMIT 10`
	}

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Private, &p.CreatedAt, &p.UpdatedAt, &p.AuthorID)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (m *PostModel) GetPosts(includePrivatePosts bool, page int, pageSize int) ([]Post, error) {
	var stmt string
	offset := (page - 1) * pageSize

	if includePrivatePosts {
		stmt = `SELECT id, title, content, private, created_at, updated_at, author_id FROM posts ORDER BY created_at DESC LIMIT ? OFFSET ?`
	} else {
		stmt = `SELECT id, title, content, private, created_at, updated_at, author_id FROM posts WHERE private = false ORDER BY created_at DESC LIMIT ? OFFSET ?`
	}

	rows, err := m.DB.Query(stmt, pageSize, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Private, &p.CreatedAt, &p.UpdatedAt, &p.AuthorID)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}
