package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/timenglesf/personal-site/internal/models"
)

var version = "1.0.0"

type application struct {
	logger *slog.Logger
	cfg    *config
	meta   *models.MetaModel
	db     *sql.DB
}

type config struct {
	port string
	db   struct {
		dsn string
	}
}

func main() {
	var cfg config
	flag.StringVar(&cfg.port, "port", os.Getenv("port"), "HTTP server port")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "./db/app.db", "SQLite3 DSN")

	flag.Parse()

	if cfg.port == "" {
		cfg.port = "4000"
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(cfg.db.dsn)
	if err != nil {
		logger.Error("Unable to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("Successfully connected to the database", "dsn", cfg.db.dsn)

	app := &application{
		logger: logger,
		cfg:    &cfg,
		meta:   &models.MetaModel{DB: db},
		db:     db,
	}

	meta := models.Meta{
		Version:     "0.0.0",
		Name:        "personal-site",
		LastUpdated: "2024-07-22",
		Description: "A simple web application",
		Author:      "Tim Engle",
		Environment: "Development",
		BuildNumber: "1",
		License:     "MIT",
	}

	fetchedMeta, err := app.meta.GetMostRecentMeta()
	if err != nil {
		logger.Error("Unable to fetch most recent meta", "error", err)
	}

	if fetchedMeta == nil || fetchedMeta.Version != meta.Version || fetchedMeta.LastUpdated != meta.LastUpdated {
		err = app.meta.InsertMeta(meta)
		if err != nil {
			logger.Error("Unable to insert meta", "error", err)
		} else {
			logger.Info("Successfully inserted meta")
		}
	}

	app.logger.Info("Successfully fetched meta", "meta", fetchedMeta)

	logger.Info("Starting the server", "port", cfg.port)
	err = http.ListenAndServe(":"+cfg.port, app.routes())
	logger.Error("Server error", "error", err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
