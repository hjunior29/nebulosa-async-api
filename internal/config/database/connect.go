package database

import (
	"log"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"github.com/hjunior29/nebulosa-async-api/internal/config"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var instance *gorm.DB

func New() error {
	dbPath := config.DATABASE_URL
	if dbPath == "" {
		dbPath = "/data/nebulosa.db"
	}

	if dir := filepath.Dir(dbPath); dir != "" && dir != "." {
		_ = os.MkdirAll(dir, 0755)
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	if sqlDB, err := db.DB(); err == nil {
		if _, err := sqlDB.Exec("PRAGMA journal_mode=WAL;"); err != nil {
			log.Println("Failed to set PRAGMA journal_mode=WAL:", err)
		}
	}

	if db != nil {
		instance = db
	}

	return nil
}

func Get() *gorm.DB {
	if instance == nil {
		log.Println("Database connection not found!")
		maxAttempts := 3
		for attempt := 0; attempt < maxAttempts; attempt++ {
			log.Println("retrying connect... attempt: ", attempt)
			err := New()
			if err != nil {
				log.Fatal(err)
			}
			if instance != nil {
				log.Println("Database connected!")
				return instance
			}
		}
	}
	return instance
}
