package config

import (
	"crypto/rsa"
	"log"
	"os"
	"strings"

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
	DATABASE_URL  string

	API_URL    string
	ORIGIN_URL string

	USERNAME string
	PASSWORD string

	PRIVATE_KEY *rsa.PrivateKey
	PUBLIC_KEY  *rsa.PublicKey
)

func init() {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")

	PORT = getEnv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	DATABASE_URL = getEnv("DATABASE_URL")
	if DATABASE_URL == "" {
		DATABASE_URL = "/data/nebulosa.db"
	}

	API_URL = getEnv("API_URL")
	if API_URL == "" {
		API_URL = "http://localhost:8080"
	}
	ORIGIN_URL = getEnv("ORIGIN_URL")
	if ORIGIN_URL == "" {
		ORIGIN_URL = "*"
	}

	USERNAME = getEnv("USERNAME")
	PASSWORD = getEnv("PASSWORD")

	privKeyStr := getEnv("PRIVATE_KEY")
	if privKeyStr != "" {
		privKeyStr = strings.ReplaceAll(privKeyStr, "\\n", "\n")
		var err error
		PRIVATE_KEY, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(privKeyStr))
		if err != nil {
			log.Println("Warning: error parsing private key:", err)
		}
	}

	pubKeyStr := getEnv("PUBLIC_KEY")
	if pubKeyStr != "" {
		pubKeyStr = strings.ReplaceAll(pubKeyStr, "\\n", "\n")
		var err error
		PUBLIC_KEY, err = jwt.ParseRSAPublicKeyFromPEM([]byte(pubKeyStr))
		if err != nil {
			log.Println("Warning: error parsing public key:", err)
		}
	}
}

func getEnv(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return ""
}
