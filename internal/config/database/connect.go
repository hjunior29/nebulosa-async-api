package database

import (
	"fmt"
	"log"

	"github.com/hjunior29/nebulosa-async-api/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var instance *gorm.DB

func New() error {
	host := config.DATABASE_HOST
	user := config.DATABASE_USER
	password := config.DATABASE_PASS
	name := config.DATABASE_NAME
	port := config.DATABASE_PORT

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, name, port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
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
