package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	PORT          = getEnv("PORT")
	DATABASE_HOST = getEnv("DATABASE_HOST")
	DATABASE_USER = getEnv("DATABASE_USER")
	DATABASE_PASS = getEnv("DATABASE_PASS")
	DATABASE_NAME = getEnv("DATABASE_NAME")
	DATABASE_PORT = getEnv("DATABASE_PORT")

	API_URL = getEnv("API_URL")
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		err = nil
		err = godotenv.Load("../.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

func getEnv(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return ""
}
