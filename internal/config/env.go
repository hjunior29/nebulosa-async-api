package config

import (
	"crypto/rsa"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var (
	PORT string

	DATABASE_HOST string
	DATABASE_USER string
	DATABASE_PASS string
	DATABASE_NAME string
	DATABASE_PORT string

	API_URL string

	USERNAME string
	PASSWORD string

	PRIVATE_KEY *rsa.PrivateKey
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

	PORT = getEnv("PORT")
	DATABASE_HOST = getEnv("DATABASE_HOST")
	DATABASE_USER = getEnv("DATABASE_USER")
	DATABASE_PASS = getEnv("DATABASE_PASS")
	DATABASE_NAME = getEnv("DATABASE_NAME")
	DATABASE_PORT = getEnv("DATABASE_PORT")

	API_URL = getEnv("API_URL")

	USERNAME = getEnv("USERNAME")
	PASSWORD = getEnv("PASSWORD")

	PRIVATE_KEY, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(getEnv("PRIVATE_KEY")))
	if err != nil {
		log.Fatal("Error parsing private key")
	}
}

func getEnv(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return ""
}
