package utils

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
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

func NotAuthorized(c *gin.Context, message string) {
	response := domain.Response{
		Status:  401,
		Message: message,
		Data:    nil,
	}
	c.AbortWithStatusJSON(401, response)
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

func VerifyJWT(tokenStr string) (bool, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return false, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return config.PUBLIC_KEY, nil
	}

	token, err := jwt.Parse(tokenStr, keyFunc)
	if err != nil {
		return false, err
	}

	if !token.Valid {
		return false, nil
	}

	return true, nil
}

func ParseScheduledAt(scheduledAt string) (time.Time, error) {
	now := time.Now()
	re := regexp.MustCompile(`^(?P<value>\d+)(?P<unit>[smhd])$`)
	matches := re.FindStringSubmatch(scheduledAt)

	if len(matches) != 3 {
		return time.Time{}, fmt.Errorf("scheduledAt format is incorrect")
	}

	value, err := strconv.Atoi(matches[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse the value in scheduledAt: %v", err)
	}

	unit := matches[2]

	switch unit {
	case "s":
		return now.Add(time.Duration(value) * time.Second), nil
	case "m":
		return now.Add(time.Duration(value) * time.Minute), nil
	case "h":
		return now.Add(time.Duration(value) * time.Hour), nil
	case "d":
		return now.Add(time.Duration(value) * 24 * time.Hour), nil
	default:
		return time.Time{}, fmt.Errorf("unknown unit in scheduledAt: %s", unit)
	}
}
