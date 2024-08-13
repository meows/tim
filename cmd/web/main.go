package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/gormstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/mattn/go-sqlite3"
	"github.com/timenglesf/personal-site/internal/models"
	"github.com/timenglesf/personal-site/ui/template"
)

const (
	sessionUserId  = "authenticatedUserID"
	sessionIsAdmin = "isAdmin"
)

var version = "0.0.0"

type application struct {
	logger            *slog.Logger
	cfg               *config
	meta              *models.MetaModel
	user              *models.UserModel
	post              *models.PostModel
	db                *gorm.DB
	sessionManager    *scs.SessionManager
	formDecoder       *form.Decoder
	pageTemplates     *template.Pages
	partialTemplates  *template.Partials
	mostRecentPost    *models.Post
	latestPublicPosts *[]models.Post
}

type config struct {
	port string
	db   struct {
		dsn string
	}
	objectStorage objectStorageConfig
}

type objectStorageConfig struct {
	objectStorageURL         string
	serveStaticObjectStorage bool
}

func main() {
	var cfg config
	flag.StringVar(&cfg.port, "port", os.Getenv("port"), "HTTP server port")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "./app.db", "SQLite3 DSN")
	flag.BoolVar(&cfg.objectStorage.serveStaticObjectStorage, "object-storage", false, "Serve static files from object storage")

	flag.Parse()

	if cfg.port == "" {
		cfg.port = "4000"
	}

	if cfg.objectStorage.serveStaticObjectStorage {
		osURL := os.Getenv("OBJECT_STORAGE_URL")
		if osURL == "" {
			log.Fatal("OBJECT_STORAGE_URL must be set when object storage is enabled")
		}
		targetFile := fmt.Sprintf("%s/static/dist/js/form-prevent.js", osURL)
		// #nosec G107
		resp, err := http.Get(targetFile)
		if err != nil {
			log.Fatal("Unable to connect to object storage")
		}
		if resp.StatusCode != http.StatusOK {
			log.Fatal("Unable to connect to object storage")
		}
		cfg.objectStorage.objectStorageURL = osURL
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := gorm.Open(sqlite.Open(cfg.db.dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&models.User{}, &models.Post{}, &models.Tag{}, &models.Meta{})
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("Successfully connected to the database", "dsn", cfg.db.dsn)

	sessionManager := scs.New()
	sessionManager.Store, err = gormstore.New(db)
	if err != nil {
		log.Fatal(err)
	}
	sessionManager.Lifetime = 24 * 7 * time.Hour

	// Initialize form decoder
	formDecoder := form.NewDecoder()

	// Initialize page templates
	pageTemplates := template.CreatePageTemplates()
	// Initialize partial templates
	partialTemplates := template.CreatePartialTemplates()

	app := &application{
		logger:           logger,
		cfg:              &cfg,
		meta:             &models.MetaModel{DB: db},
		user:             &models.UserModel{DB: db},
		post:             &models.PostModel{DB: db},
		db:               db,
		sessionManager:   sessionManager,
		formDecoder:      formDecoder,
		pageTemplates:    pageTemplates,
		partialTemplates: partialTemplates,
	}

	// Reset mostRecentPublicPost & latestPublicPosts app field
	if err := app.UpdatePostsOnAppStruct(); err != nil {
		app.logger.Error("Unable to update posts on app struct", "error", err)
	}

	meta := models.Meta{
		Version:     version,
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
		app.logger.Error("Unable to fetch most recent meta", "error", err)
	}

	if fetchedMeta == nil || fetchedMeta.Version != meta.Version || fetchedMeta.LastUpdated != meta.LastUpdated {
		err = app.meta.InsertMeta(meta)
		if err != nil {
			logger.Error("Unable to insert meta", "error", err)
		} else {
			logger.Info("Successfully inserted meta")
		}
	}

	srv := &http.Server{
		Addr:         ":" + cfg.port,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.logger.Info("Successfully fetched meta", "meta", fetchedMeta)

	logger.Info("Starting the server", "port", cfg.port)
	err = srv.ListenAndServe()
	logger.Error("Server error", "error", err.Error())
	os.Exit(1)
}
