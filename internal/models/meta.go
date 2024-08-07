package models

import (
	"gorm.io/gorm"
)

type Meta struct {
	gorm.Model
	Version     string `gorm:"not null"`
	Name        string `gorm:"not null"`
	LastUpdated string `gorm:"not null"`
	Description string
	Author      string
	Environment string
	BuildNumber string
	License     string
}

type MetaModel struct {
	DB *gorm.DB
}

func (m *MetaModel) InsertMeta(md Meta) error {
	result := m.DB.Create(&md)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (m *MetaModel) GetMostRecentMeta() (*Meta, error) {
	var meta Meta
	result := m.DB.Order("id desc").First(&meta)
	if result.Error != nil {
		return nil, result.Error
	}
	return &meta, nil
}

// type Meta struct {
// 	ID          int
// 	Version     string
// 	Name        string
// 	LastUpdated string
// 	Description string
// 	Author      string
// 	Environment string
// 	BuildNumber string
// 	License     string
// }
// type MetaModel struct {
// 	DB *sql.DB
// }
//
// type MetalModelInterface interface{}
//
// func (m *MetaModel) InsertMeta(md Meta) error {
// 	query := `
//     INSERT INTO meta (version, name, last_updated, description, author, environment, build_number, license)
//     VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
//
// 	_, err := m.DB.Exec(query, md.Version, md.Name, md.LastUpdated, md.Description, md.Author, md.Environment, md.BuildNumber, md.License)
// 	return err
// }
//
// func (m *MetaModel) GetMostRecentMeta() (*Meta, error) {
// 	var meta Meta
// 	query := `SELECT id, version, name, last_updated, description, author, environment, build_number, license FROM meta ORDER BY id DESC LIMIT 1`
// 	row := m.DB.QueryRow(query)
// 	err := row.Scan(&meta.ID, &meta.Version, &meta.Name, &meta.LastUpdated, &meta.Description, &meta.Author, &meta.Environment, &meta.BuildNumber, &meta.License)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &meta, nil
// }
