package utils

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hjunior29/nebulosa-async-api/internal/config"
	"github.com/hjunior29/nebulosa-async-api/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func SuccessResponse(c *gin.Context, status int, message string, data interface{}) {
	response := domain.Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
	c.JSON(status, response)
}

func ErrorResponse(c *gin.Context, status int, message string, err error) {
	if err != nil {
		log.Printf("[ERROR] %s\n", err.Error())

		if message == "" {
			message = err.Error()
		}
	}

	response := domain.Response{
		Status:  status,
		Message: message,
		Data:    nil,
	}
	c.JSON(status, response)

}

func GetId(c *gin.Context) uuid.UUID {
	idStr := c.Param("id")
	if idStr == "" {
		ErrorResponse(c, 400, "Missing id", nil)
		return uuid.UUID{}
	}

	return uuid.MustParse(idStr)
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GenerateJWT(id, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"id":       id,
		"username": username,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenStr, err := token.SignedString(config.PRIVATE_KEY)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
