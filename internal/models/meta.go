package models

import "database/sql"

type Meta struct {
	ID          int
	Version     string
	Name        string
	LastUpdated string
	Description string
	Author      string
	Environment string
	BuildNumber string
	License     string
}

type MetaModel struct {
	DB *sql.DB
}

type MetalModelInterface interface{}

func (m *MetaModel) InsertMeta(md Meta) error {
	query := `
    INSERT INTO meta (version, name, last_updated, description, author, environment, build_number, license) 
    VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := m.DB.Exec(query, md.Version, md.Name, md.LastUpdated, md.Description, md.Author, md.Environment, md.BuildNumber, md.License)
	return err
}

func (m *MetaModel) GetMostRecentMeta() (*Meta, error) {
	var meta Meta
	query := `SELECT id, version, name, last_updated, description, author, environment, build_number, license FROM meta ORDER BY id DESC LIMIT 1`
	row := m.DB.QueryRow(query)
	err := row.Scan(&meta.ID, &meta.Version, &meta.Name, &meta.LastUpdated, &meta.Description, &meta.Author, &meta.Environment, &meta.BuildNumber, &meta.License)
	if err != nil {
		return nil, err
	}
	return &meta, nil
}
