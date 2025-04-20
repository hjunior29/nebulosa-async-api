package database

import (
	"log"

	"github.com/hjunior29/nebulosa-async-api/internal/config"
	"github.com/hjunior29/nebulosa-async-api/internal/domain"
	"github.com/hjunior29/nebulosa-async-api/internal/utils"
	"gorm.io/gorm"
)

func Migrate(models ...interface{}) error {
	err := Get().AutoMigrate(models...)
	if err != nil {
		return err
	}

	log.Println("Database migration completed successfully.")
	return nil
}

func Seed() error {
	var user domain.User
	repository := NewRepository(&user, nil)

	username := config.USERNAME
	password := config.PASSWORD

	if username == "" || password == "" {
		log.Panic("Username or password not provided. Skipping seeding.")
		return nil
	}

	err := repository.FindAllWhere(map[string]interface{}{"username": username})
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Println("Error finding user:", err)
			return err
		}
		log.Println("User not found. Proceeding with seeding.")
	}

	if user.Username != "" {
		log.Println("User already exists. Skipping seeding.")
		return nil
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Panic("Error hashing password:", err)
	}

	user.Username = username
	user.HashedPassword = hashedPassword
	err = repository.Create()
	if err != nil {
		log.Panic("Error creating user:", err)
	}

	log.Println("Database seeding completed successfully.")
	return nil
}
