package task

import (
	"github.com/gin-gonic/gin"
	"github.com/hjunior29/nebulosa-async-api/internal/config/database"
	"github.com/hjunior29/nebulosa-async-api/internal/domain"
	"github.com/hjunior29/nebulosa-async-api/internal/utils"
)

func Create(c *gin.Context) {
	var task domain.Task
	repository := database.NewRepository(&task, c)

	err := c.ShouldBindJSON(&task)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid request", err)
		return
	}

	if task.Endpoint == "" || task.Type == "" || task.Payload == nil || task.MaxRetries == 0 {
		utils.ErrorResponse(c, 400, "Missing required fields", nil)
		return
	}

	task.Status = domain.StatusPending

	err = repository.Create()
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to create task", err)
		return
	}

	utils.SuccessResponse(c, 201, "Task created successfully", task)

}
