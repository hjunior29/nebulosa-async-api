package task

import (
	"github.com/gin-gonic/gin"
	"github.com/hjunior29/nebulosa-async-api/internal/config/database"
	"github.com/hjunior29/nebulosa-async-api/internal/domain"
	"github.com/hjunior29/nebulosa-async-api/internal/utils"
)

func Delete(c *gin.Context) {
	var task domain.Task
	repository := database.NewRepository(&task, c)

	id := utils.GetId(c)

	err := repository.GetById(id)
	if err != nil {
		utils.ErrorResponse(c, 404, "Task not found", err)
		return
	}

	err = repository.Delete(id)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to delete task", err)
		return
	}

	utils.SuccessResponse(c, 200, "Task found", task)
}
