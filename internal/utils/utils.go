package utils

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hjunior29/nebulosa-async-api/internal/domain"
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
