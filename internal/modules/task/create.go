package task

import (
	"fmt"
	"log"

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

	var requestBody map[string]interface{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request body", err)
		return
	}

	if interval, ok := requestBody["scheduledAtInterval"].(string); ok && interval != "" {
		futureTime, parseErr := utils.ParseScheduledAt(interval)
		if parseErr != nil {
			utils.ErrorResponse(c, 400, "Invalid scheduledAtInterval format", parseErr)
			return
		}
		task.ScheduledAtTime = futureTime
	}

	task.Status = domain.StatusPending

	err = repository.Create()
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to create task", err)
		return
	}

	sqlDB, dbErr := database.Get().DB()
	if dbErr != nil {
		log.Println("Failed to get sql.DB:", dbErr)
	} else {
		_, notifyErr := sqlDB.Exec(fmt.Sprintf("NOTIFY new_task, '%s'", task.ID.String()))
		if notifyErr != nil {
			log.Println("Failed to notify new_task:", notifyErr)
		}
	}

	utils.SuccessResponse(c, 201, "Task created successfully", task)

}
