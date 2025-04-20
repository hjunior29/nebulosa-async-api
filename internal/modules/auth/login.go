package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/hjunior29/nebulosa-async-api/internal/config/database"
	"github.com/hjunior29/nebulosa-async-api/internal/domain"
	"github.com/hjunior29/nebulosa-async-api/internal/utils"
)

func Login(c *gin.Context) {
	var userRequest domain.UserRequest
	var user domain.User
	repository := database.NewRepository(&user, c)

	err := c.ShouldBindJSON(&userRequest)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid request", err)
		return
	}

	if userRequest.Username == "" || userRequest.Password == "" {
		utils.ErrorResponse(c, 400, "Username and password are required", nil)
		return
	}

	err = repository.FindAllWhere(map[string]interface{}{
		"username": userRequest.Username,
	})
	if err != nil {
		utils.ErrorResponse(c, 500, "Error finding user", err)
		return
	}

	isValid := utils.VerifyPassword(user.HashedPassword, userRequest.Password)
	if !isValid {
		utils.ErrorResponse(c, 401, "Invalid username or password", nil)
		return
	}

	token, err := utils.GenerateJWT(user.ID.String(), user.Username)
	if err != nil {
		utils.ErrorResponse(c, 500, "Error generating token", err)
		return
	}

	utils.SuccessResponse(c, 200, "Login successful", map[string]interface{}{"token": token})
}
